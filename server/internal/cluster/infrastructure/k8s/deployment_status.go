package k8s

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"starliner.app/internal/cluster/domain/port"
)

const deploymentStatusPollInterval = 2 * time.Second

// Visible between snapshots; each line is sent as its own SSE event.
const deploymentStatusSnapshotDelimiter = "────────────────────────────────────────"

var (
	pulledDurationRe = regexp.MustCompile(`in ([0-9.]+s)`)
	pulledSizeRe     = regexp.MustCompile(`Image size:\s*([0-9]+)\s*bytes`)
)

type DeploymentStatus struct{}

func NewDeploymentStatus() port.DeploymentStatus {
	return &DeploymentStatus{}
}

func (d *DeploymentStatus) StreamDeploymentStatusLogs(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	commitHash string,
) (io.ReadCloser, error) {
	reader, writer := io.Pipe()

	go func() {
		defer func() {
			_ = writer.Close()
		}()

		firstSnapshot := true

		for {
			select {
			case <-ctx.Done():
				_ = writer.CloseWithError(ctx.Err())
				return
			default:
			}

			report, done, err := d.collectReport(ctx, namespace, releaseName, kubeconfigBase64, commitHash)
			if err != nil {
				if ctx.Err() != nil {
					_ = writer.CloseWithError(ctx.Err())
					return
				}
				if _, werr := fmt.Fprintf(writer, "Error: %v\n\n", err); werr != nil {
					return
				}
				time.Sleep(deploymentStatusPollInterval)
				continue
			}

			if !firstSnapshot {
				if _, err := fmt.Fprintf(writer, "%s\n", deploymentStatusSnapshotDelimiter); err != nil {
					return
				}
			}
			firstSnapshot = false

			if _, err := io.WriteString(writer, report); err != nil {
				return
			}

			if done {
				return
			}

			select {
			case <-ctx.Done():
				_ = writer.CloseWithError(ctx.Err())
				return
			case <-time.After(deploymentStatusPollInterval):
			}
		}
	}()

	return reader, nil
}

type podCounts struct {
	running     int
	starting    int
	terminating int
	failed      int
}

func (c podCounts) String() string {
	return fmt.Sprintf(
		"%d running, %d starting, %d terminating, %d failed",
		c.running,
		c.starting,
		c.terminating,
		c.failed,
	)
}

type podProgress struct {
	pod    corev1.Pod
	status string
	events []podEvent
}

type podEvent struct {
	kind     string
	detail   string
	extra    []string
	complete bool
}

func (d *DeploymentStatus) collectReport(
	ctx context.Context,
	namespace string,
	releaseName string,
	kubeconfigBase64 string,
	commitHash string,
) (string, bool, error) {
	client, err := newKubernetesClient(kubeconfigBase64)
	if err != nil {
		return "", false, err
	}

	deploy, err := client.AppsV1().Deployments(namespace).Get(ctx, releaseName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return d.collectStatefulSetReport(ctx, client, namespace, releaseName)
		}
		return "", false, fmt.Errorf("get deployment: %w", err)
	}

	labelSelector := metav1.FormatLabelSelector(deploy.Spec.Selector)

	rsList, err := client.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", false, fmt.Errorf("list replicasets: %w", err)
	}

	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", false, fmt.Errorf("list pods: %w", err)
	}

	newRS, oldRS := splitReplicaSets(deploy, rsList.Items)
	versionLabel := commitHash
	if versionLabel == "" {
		versionLabel = versionFromReplicaSet(newRS)
	}
	if versionLabel == "" {
		versionLabel = "unknown"
	}

	prevCounts := countPodsForReplicaSets(pods.Items, oldRS)
	newCounts := countNewVersionPods(pods.Items, newRS, oldRS)

	startingPods := collectStartingPods(ctx, client, namespace, pods.Items, newRS, oldRS)

	failed := newCounts.failed > 0 || deploymentFailed(deploy)
	success := rolloutHealthy(deploy, newCounts, prevCounts)
	done := failed || success

	var b strings.Builder
	b.WriteString("🚀 Deployment Status Report\n\n")

	var statusLine string
	switch {
	case failed:
		statusLine = "Application deployment for commit %s has failed.\n\n"
	case success:
		statusLine = "Application deployment for commit %s is complete.\n\n"
	default:
		statusLine = "Application deployment for commit %s is in progress.\n\n"
	}
	if _, err := fmt.Fprintf(&b, statusLine, versionLabel); err != nil {
		return "", false, err
	}

	b.WriteString("📦 Previous Version\n")
	if _, err := fmt.Fprintf(&b, "└─ Pods: %s\n\n", prevCounts); err != nil {
		return "", false, err
	}
	if _, err := fmt.Fprintf(&b, "📦 New Version (%s)\n", versionLabel); err != nil {
		return "", false, err
	}
	if _, err := fmt.Fprintf(&b, "└─ Pods: %s\n", newCounts); err != nil {
		return "", false, err
	}

	for _, p := range startingPods {
		if err := formatPodProgress(&b, p); err != nil {
			return "", false, err
		}
	}

	if !success && !failed && (newCounts.starting > 0 || newCounts.running < desiredReplicas(deploy)) {
		b.WriteString("\n⏳ Waiting for pod readiness checks to pass...\n")
	}

	return b.String(), done, nil
}

func (d *DeploymentStatus) collectStatefulSetReport(
	ctx context.Context,
	client *kubernetes.Clientset,
	namespace string,
	releaseName string,
) (string, bool, error) {
	stsName := releaseName + "-db"
	sts, err := client.AppsV1().StatefulSets(namespace).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Sprintf(
				"🚀 Deployment Status Report\n\nWaiting for database %s to appear in the cluster...\n",
				stsName,
			), false, nil
		}
		return "", false, fmt.Errorf("get statefulset: %w", err)
	}

	labelSelector := metav1.FormatLabelSelector(sts.Spec.Selector)
	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", false, fmt.Errorf("list pods: %w", err)
	}

	counts := countAllPods(pods.Items)
	versionLabel := versionFromStatefulSet(sts)
	startingPods := collectStartingPods(ctx, client, namespace, pods.Items, nil, nil)

	desired := desiredStatefulSetReplicas(sts)
	failed := counts.failed > 0
	success := !failed && counts.running >= desired && counts.starting == 0 && counts.terminating == 0
	done := failed || success

	var b strings.Builder
	b.WriteString("🚀 Deployment Status Report\n\n")

	var statusLine string
	switch {
	case failed:
		statusLine = "Application deployment for commit %s has failed.\n\n"
	case success:
		statusLine = "Application deployment for commit %s is complete.\n\n"
	default:
		statusLine = "Application deployment for commit %s is in progress.\n\n"
	}
	if _, err := fmt.Fprintf(&b, statusLine, versionLabel); err != nil {
		return "", false, err
	}

	if _, err := fmt.Fprintf(&b, "📦 Pods\n└─ %s\n", counts); err != nil {
		return "", false, err
	}

	for _, p := range startingPods {
		if err := formatPodProgress(&b, p); err != nil {
			return "", false, err
		}
	}

	if !success && !failed && (counts.starting > 0 || counts.running < desired) {
		b.WriteString("\n⏳ Waiting for pod readiness checks to pass...\n")
	}

	return b.String(), done, nil
}

func desiredStatefulSetReplicas(sts *appsv1.StatefulSet) int {
	if sts.Spec.Replicas == nil {
		return 1
	}
	return int(*sts.Spec.Replicas)
}

func versionFromStatefulSet(sts *appsv1.StatefulSet) string {
	if len(sts.Spec.Template.Spec.Containers) == 0 {
		return "unknown"
	}
	image := sts.Spec.Template.Spec.Containers[0].Image
	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return image
	}
	return parts[len(parts)-1]
}

func desiredReplicas(deploy *appsv1.Deployment) int {
	if deploy.Spec.Replicas == nil {
		return 1
	}
	return int(*deploy.Spec.Replicas)
}

func splitReplicaSets(deploy *appsv1.Deployment, sets []appsv1.ReplicaSet) (*appsv1.ReplicaSet, []*appsv1.ReplicaSet) {
	var owned []appsv1.ReplicaSet
	for i := range sets {
		rs := sets[i]
		if !metav1.IsControlledBy(&rs, deploy) {
			continue
		}
		owned = append(owned, rs)
	}

	if len(owned) == 0 {
		return nil, nil
	}

	sort.Slice(owned, func(i, j int) bool {
		return replicaSetRevision(&owned[i]) > replicaSetRevision(&owned[j])
	})

	newRS := &owned[0]
	var old []*appsv1.ReplicaSet
	for i := 1; i < len(owned); i++ {
		old = append(old, &owned[i])
	}
	return newRS, old
}

func replicaSetRevision(rs *appsv1.ReplicaSet) int {
	if rs == nil {
		return 0
	}
	rev, _ := strconv.Atoi(rs.Annotations["deployment.kubernetes.io/revision"])
	return rev
}

func versionFromReplicaSet(rs *appsv1.ReplicaSet) string {
	if rs == nil {
		return ""
	}
	if len(rs.Spec.Template.Spec.Containers) == 0 {
		return ""
	}
	image := rs.Spec.Template.Spec.Containers[0].Image
	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return image
	}
	tag := parts[len(parts)-1]
	if idx := strings.LastIndex(tag, "-"); idx >= 0 {
		candidate := tag[idx+1:]
		if len(candidate) >= 7 {
			return candidate
		}
	}
	return tag
}

func podBelongsToReplicaSet(pod corev1.Pod, rs *appsv1.ReplicaSet) bool {
	if rs == nil {
		return false
	}
	podHash := pod.Labels["pod-template-hash"]
	rsHash := rs.Labels["pod-template-hash"]
	return podHash != "" && podHash == rsHash
}

func countPodsForReplicaSets(pods []corev1.Pod, sets []*appsv1.ReplicaSet) podCounts {
	var counts podCounts
	for _, pod := range pods {
		if !podMatchesAnyReplicaSet(pod, sets) {
			continue
		}
		switch classifyPod(&pod) {
		case "running":
			counts.running++
		case "starting":
			counts.starting++
		case "terminating":
			counts.terminating++
		case "failed":
			counts.failed++
		}
	}
	return counts
}

func countNewVersionPods(pods []corev1.Pod, newRS *appsv1.ReplicaSet, oldRS []*appsv1.ReplicaSet) podCounts {
	if newRS != nil {
		return countPodsForReplicaSets(pods, []*appsv1.ReplicaSet{newRS})
	}
	if len(oldRS) == 0 {
		return countAllPods(pods)
	}
	return podCounts{}
}

func countAllPods(pods []corev1.Pod) podCounts {
	var counts podCounts
	for _, pod := range pods {
		switch classifyPod(&pod) {
		case "running":
			counts.running++
		case "starting":
			counts.starting++
		case "terminating":
			counts.terminating++
		case "failed":
			counts.failed++
		}
	}
	return counts
}

func podMatchesAnyReplicaSet(pod corev1.Pod, sets []*appsv1.ReplicaSet) bool {
	if len(sets) == 0 {
		return false
	}
	for _, rs := range sets {
		if rs != nil && podBelongsToReplicaSet(pod, rs) {
			return true
		}
	}
	return false
}

func classifyPod(pod *corev1.Pod) string {
	if pod.DeletionTimestamp != nil {
		return "terminating"
	}

	if pod.Status.Phase == corev1.PodFailed {
		return "failed"
	}

	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil {
			switch cs.State.Waiting.Reason {
			case "CrashLoopBackOff", "ImagePullBackOff", "ErrImagePull", "Error":
				return "failed"
			}
		}
	}

	if pod.Status.Phase == corev1.PodRunning {
		for _, cs := range pod.Status.ContainerStatuses {
			if !cs.Ready {
				return "starting"
			}
		}
		return "running"
	}

	return "starting"
}

func collectStartingPods(
	ctx context.Context,
	client *kubernetes.Clientset,
	namespace string,
	pods []corev1.Pod,
	newRS *appsv1.ReplicaSet,
	oldRS []*appsv1.ReplicaSet,
) []podProgress {
	var result []podProgress
	for _, pod := range pods {
		if newRS != nil {
			if !podBelongsToReplicaSet(pod, newRS) {
				continue
			}
		} else if podMatchesAnyReplicaSet(pod, oldRS) {
			continue
		}

		status := classifyPod(&pod)
		if status != "starting" && status != "failed" {
			continue
		}
		result = append(result, podProgress{
			pod:    pod,
			status: strings.ToUpper(status),
			events: describePodEvents(ctx, client, namespace, pod),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].pod.Name < result[j].pod.Name
	})
	return result
}

func describePodEvents(
	ctx context.Context,
	client *kubernetes.Clientset,
	namespace string,
	pod corev1.Pod,
) []podEvent {
	events, err := client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=Pod", pod.Name),
	})
	if err != nil {
		return nil
	}

	sort.Slice(events.Items, func(i, j int) bool {
		return eventTime(events.Items[i]).Before(eventTime(events.Items[j]))
	})

	var progress []podEvent
	seen := map[string]bool{}

	for _, event := range events.Items {
		switch event.Reason {
		case "Pulled":
			if seen["Pulled"] {
				continue
			}
			seen["Pulled"] = true
			progress = append(progress, parsePulledEvent(event))
		case "Created":
			if seen["Created"] {
				continue
			}
			seen["Created"] = true
			container := containerFromEvent(event, pod)
			progress = append(progress, podEvent{
				kind:     "Container created",
				detail:   container,
				complete: true,
			})
		case "Started":
			if seen["Started"] {
				continue
			}
			seen["Started"] = true
			container := containerFromEvent(event, pod)
			progress = append(progress, podEvent{
				kind:     "Container started",
				detail:   container,
				complete: true,
			})
		}
	}

	if len(progress) == 0 {
		progress = progressFromContainerStatus(pod)
	}

	return progress
}

func eventTime(event corev1.Event) time.Time {
	if !event.LastTimestamp.IsZero() {
		return event.LastTimestamp.Time
	}
	if !event.EventTime.IsZero() {
		return event.EventTime.Time
	}
	return event.CreationTimestamp.Time
}

func parsePulledEvent(event corev1.Event) podEvent {
	image := strings.TrimSpace(event.Message)
	if idx := strings.Index(image, `"`); idx >= 0 {
		if end := strings.Index(image[idx+1:], `"`); end >= 0 {
			image = image[idx+1 : idx+1+end]
		}
	}

	var extra []string
	if match := pulledDurationRe.FindStringSubmatch(event.Message); len(match) > 1 {
		extra = append(extra, fmt.Sprintf("Duration: %s", match[1]))
	}
	if match := pulledSizeRe.FindStringSubmatch(event.Message); len(match) > 1 {
		sizeBytes, _ := strconv.ParseInt(match[1], 10, 64)
		extra = append(extra, fmt.Sprintf("Size: %s", formatBytesSI(sizeBytes)))
	}

	return podEvent{
		kind:     "Image pulled",
		detail:   image,
		extra:    extra,
		complete: true,
	}
}

func formatBytesSI(bytes int64) string {
	return fmt.Sprintf("%.1f MB", float64(bytes)/1_000_000)
}

func containerFromEvent(event corev1.Event, pod corev1.Pod) string {
	if strings.Contains(event.Message, "container ") {
		parts := strings.Fields(event.Message)
		for i, part := range parts {
			if part == "container" && i+1 < len(parts) {
				name := strings.Trim(parts[i+1], `":`)
				if name != "" {
					return name
				}
			}
		}
	}

	if len(pod.Spec.Containers) > 0 {
		return pod.Spec.Containers[0].Name
	}
	return pod.Name
}

func progressFromContainerStatus(pod corev1.Pod) []podEvent {
	var progress []podEvent
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil {
			switch cs.State.Waiting.Reason {
			case "ContainerCreating":
				progress = append(progress, podEvent{
					kind:     "Container creating",
					detail:   cs.Name,
					complete: false,
				})
			case "PodInitializing":
				progress = append(progress, podEvent{
					kind:     "Pod initializing",
					detail:   cs.Name,
					complete: false,
				})
			}
		}
		if cs.State.Running != nil && !cs.Ready {
			progress = append(progress, podEvent{
				kind:     "Container running",
				detail:   cs.Name,
				complete: false,
			})
		}
	}
	return progress
}

func formatPodProgress(b *strings.Builder, p podProgress) error {
	b.WriteString("\n\n")
	if _, err := fmt.Fprintf(b, "   └─ Pod %s\n", p.pod.Name); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(b, "      Status: %s\n", p.status); err != nil {
		return err
	}

	for _, event := range p.events {
		marker := "✓"
		if !event.complete {
			marker = "…"
		}
		if _, err := fmt.Fprintf(b, "\n      %s %s\n", marker, event.kind); err != nil {
			return err
		}
		if event.detail != "" {
			if _, err := fmt.Fprintf(b, "        %s\n", event.detail); err != nil {
				return err
			}
		}
		for _, line := range event.extra {
			if _, err := fmt.Fprintf(b, "        %s\n", line); err != nil {
				return err
			}
		}
	}
	return nil
}

func deploymentFailed(deploy *appsv1.Deployment) bool {
	for _, c := range deploy.Status.Conditions {
		if c.Type == appsv1.DeploymentProgressing &&
			c.Status == corev1.ConditionFalse &&
			c.Reason == "ProgressDeadlineExceeded" {
			return true
		}
	}
	return false
}

func rolloutHealthy(deploy *appsv1.Deployment, newCounts, prevCounts podCounts) bool {
	desired := desiredReplicas(deploy)
	if newCounts.failed > 0 {
		return false
	}
	if newCounts.running < desired {
		return false
	}
	if newCounts.starting > 0 {
		return false
	}
	if prevCounts.running > 0 || prevCounts.starting > 0 {
		return false
	}

	replicas := deploy.Spec.Replicas
	if replicas == nil {
		return deploy.Status.AvailableReplicas >= int32(desired)
	}

	return deploy.Status.UpdatedReplicas >= *replicas &&
		deploy.Status.AvailableReplicas >= *replicas &&
		deploy.Status.ReadyReplicas >= *replicas
}

package ansible

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/provisioner/domain/port"
)

type Install struct {
	crypto       corePort.Crypto
	logPublisher port.LogPublisher
}

func NewInstall(
	crypto corePort.Crypto,
	logPublisher port.LogPublisher,
) port.Install {
	return &Install{
		crypto:       crypto,
		logPublisher: logPublisher,
	}
}

func (i *Install) InstallK3s(clusterId int64, ip string, privateKey []byte) (kubeconfig string, logs string, err error) {
	var (
		logBuf strings.Builder
		mu     sync.Mutex
	)
	appendLog := func(line string) {
		mu.Lock()
		logBuf.WriteString(line)
		mu.Unlock()

		if i.logPublisher != nil {
			if err := i.logPublisher.PublishLogChunk(clusterId, []byte(line)); err != nil {
				fmt.Printf("failed to publish log chunk: %v", err)
			}
		}
	}

	defer func() {
		logs = logBuf.String()
	}()

	tmpPlaybook, err := os.CreateTemp("", "ansible-*.yml")
	if err != nil {
		return "", "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPlaybook.Name())

	_, err = tmpPlaybook.Write([]byte(K3sPlaybook))
	if err != nil {
		return "", "", err
	}
	err = tmpPlaybook.Close()
	if err != nil {
		return "", "", err
	}

	tmpPrivateKey, err := os.CreateTemp("", "private-key-*.pem")
	if err != nil {
		return "", "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPrivateKey.Name())

	pemBytes, err := i.crypto.EncodePrivateKeyToPEM(privateKey)
	if err != nil {
		return "", "", err
	}
	_, err = tmpPrivateKey.Write(pemBytes)
	if err != nil {
		return "", "", err
	}
	err = tmpPrivateKey.Close()
	if err != nil {
		return "", "", err
	}

	tmpKubeconfig, err := os.CreateTemp("", "kubeconfig-*.yaml")
	if err != nil {
		return "", "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpKubeconfig.Name())
	if err := tmpKubeconfig.Close(); err != nil {
		return "", "", err
	}

	args := []string{
		tmpPlaybook.Name(),
		"-i", fmt.Sprintf("%s,", ip),
		"-u", "root",
		"--private-key", tmpPrivateKey.Name(),
		"-e", "ansible_ssh_common_args='-o StrictHostKeyChecking=no'",
		"-e", fmt.Sprintf("target_host=%s", ip),
		"-e", fmt.Sprintf("kubeconfig_dest=%s", tmpKubeconfig.Name()),
	}

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Env = append(
		os.Environ(),
		"ANSIBLE_DEPRECATION_WARNINGS=False",
		"ANSIBLE_FORCE_COLOR=0",
		"PYTHONUNBUFFERED=1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	var wg sync.WaitGroup

	readPipe := func(r io.Reader) {
		defer wg.Done()

		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

		for scanner.Scan() {
			appendLog(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			appendLog(fmt.Sprintf("error reading ansible output: %v\n", err))
		}
	}

	wg.Add(2)
	go readPipe(stdout)
	go readPipe(stderr)

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	kubeconfigBytes, err := os.ReadFile(tmpKubeconfig.Name())
	if err != nil {
		return "", "", fmt.Errorf("failed to read kubeconfig from %s: %w", tmpKubeconfig.Name(), err)
	}

	kubeconfig = strings.ReplaceAll(
		string(kubeconfigBytes),
		"https://127.0.0.1:",
		fmt.Sprintf("https://%s:", ip),
	)

	return kubeconfig, "", nil
}

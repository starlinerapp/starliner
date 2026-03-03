import React from "react";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useParams } from "react-router";
import DeployIngressForm, {
  type IngressFormInput,
} from "~/components/organisms/forms/DeployIngressForm";
import { toSlug } from "~/utils/slug";

export default function UpdateIngressForm() {
  const { id: projectId, deploymentId } = useParams<{
    id: string;
    deploymentId: string;
  }>();

  const { environment: currentEnvironment, clusterId } = useEnvironment();

  const trpc = useTRPC();
  const { data: environmentDeployments, isLoading } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: currentEnvironment?.id },
      { enabled: !!currentEnvironment },
    ),
  );

  const updateIngressMutation = useMutation(
    trpc.deployment.updateIngress.mutationOptions(),
  );

  const { data: clusterData } = useQuery(
    trpc.cluster.getCluster.queryOptions(
      { id: clusterId! },
      { enabled: !!clusterId },
    ),
  );

  const { data: projectData } = useQuery(
    trpc.project.getProject.queryOptions({
      id: Number(projectId),
    }),
  );

  const projectNameSlug = toSlug(projectData?.name ?? "");

  const onSubmit = async (data: IngressFormInput) => {
    await updateIngressMutation.mutateAsync({
      id: currentEnvironment.id,
      deploymentId: Number(deploymentId),
      ingressHosts: data.hosts.map((h) => ({
        host: h.name + `.${projectNameSlug}.${clusterData?.ipv4Address}.nip.io`,
        paths: h.paths.map((p) => ({
          path: p.path,
          pathType: p.pathType as "Prefix" | "Exact",
          serviceName: p.service,
        })),
      })),
    });
  };

  const ingressDeployment = environmentDeployments?.ingresses.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployIngressForm
          key={deploymentId}
          onSubmit={onSubmit}
          defaultValues={{
            hosts: (ingressDeployment?.hosts ?? []).map((h) => ({
              name: h.host.split(".")[0] ?? "",
              paths: (h.paths ?? []).map((p) => ({
                path: p.path ?? "",
                pathType:
                  p.pathType === "Prefix" || p.pathType === "Exact"
                    ? p.pathType
                    : "",
                service: p.serviceName ?? "",
              })),
            })),
          }}
        />
      )}
    </>
  );
}

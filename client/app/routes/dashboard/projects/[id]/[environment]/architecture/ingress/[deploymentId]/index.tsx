import React from "react";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useParams } from "react-router";
import DeployIngressForm from "~/components/organisms/forms/DeployIngressForm";

export default function UpdateIngressForm() {
  const { deploymentId } = useParams<{ deploymentId: string }>();

  const { environment: currentEnvironment } = useEnvironment();

  const trpc = useTRPC();
  const { data: environmentDeployments, isLoading } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: currentEnvironment?.id },
      { enabled: !!currentEnvironment },
    ),
  );

  const ingressDeployment = environmentDeployments?.ingresses.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployIngressForm
          key={deploymentId}
          defaultValues={{
            hosts: (ingressDeployment?.hosts ?? []).map((h) => ({
              name: h.host ?? "",
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

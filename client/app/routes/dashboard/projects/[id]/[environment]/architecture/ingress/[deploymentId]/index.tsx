import React from "react";
import { useNavigate, useLoaderData, useParams } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import DeployIngressForm, {
  type IngressFormInput,
} from "~/components/organisms/forms/DeployIngressForm";

export function loader() {
  return {
    deploymentEnvironment: process.env.ENVIRONMENT ?? "",
  };
}

export default function UpdateIngressForm() {
  const { deploymentEnvironment } = useLoaderData<typeof loader>();
  const { slug, id, environment, deploymentId } = useParams<{
    slug: string;
    id: string;
    environment: string;
    deploymentId: string;
  }>();
  const navigate = useNavigate();

  const { environment: currentEnvironment } = useEnvironment();
  const queryClient = useQueryClient();

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

  const onSubmit = async (data: IngressFormInput) => {
    await updateIngressMutation.mutateAsync(
      {
        id: currentEnvironment.id,
        deploymentId: Number(deploymentId),
        ingressHosts: data.hosts.map((h) => ({
          host: h.name,
          paths: h.paths.map((p) => ({
            path: p.path,
            pathType: p.pathType as "Prefix" | "Exact",
            serviceName: p.service,
          })),
        })),
      },
      {
        onSuccess: (result) => {
          queryClient.invalidateQueries({
            queryKey: trpc.environment.getEnvironmentBuilds.queryKey({
              id: currentEnvironment.id,
            }),
          });
          queryClient.invalidateQueries({
            queryKey: trpc.environment.getEnvironmentDeployments.queryKey({
              id: currentEnvironment.id,
            }),
          });
          if (result?.deploymentId && slug && id && environment) {
            navigate(
              `/${slug}/projects/${id}/${environment}/architecture/ingress/${result.deploymentId}`,
              { replace: true },
            );
          }
        },
      },
    );
  };

  const ingressDeployment = environmentDeployments?.ingresses.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployIngressForm
          deploymentEnvironment={deploymentEnvironment}
          key={deploymentId}
          onSubmit={onSubmit}
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

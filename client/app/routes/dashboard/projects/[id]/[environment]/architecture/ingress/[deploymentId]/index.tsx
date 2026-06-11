import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useLoaderData, useNavigate, useParams } from "react-router";
import DeployIngressForm, {
  type IngressFormInput,
} from "~/components/organisms/forms/DeployIngressForm";
import { serverEnv } from "~/env.server";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";

export function loader() {
  return {
    deploymentEnvironment: serverEnv.ENVIRONMENT,
    deploymentDomain: serverEnv.DEPLOYMENT_DOMAIN,
  };
}

export default function UpdateIngressForm() {
  const { deploymentEnvironment, deploymentDomain } =
    useLoaderData<typeof loader>();
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
          prefix: h.name,
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
          deploymentDomain={deploymentDomain}
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

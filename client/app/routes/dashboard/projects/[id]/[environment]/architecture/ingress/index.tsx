import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useLoaderData } from "react-router";
import DeployIngressForm, {
  type IngressFormInput,
} from "~/components/organisms/forms/DeployIngressForm";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";

export function loader() {
  return {
    deploymentEnvironment: process.env.ENVIRONMENT ?? "",
  };
}

export default function Index() {
  const { deploymentEnvironment } = useLoaderData<typeof loader>();
  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const createIngressMutation = useMutation(
    trpc.deployment.deployIngress.mutationOptions(),
  );
  const { environment: currentEnvironment } = useEnvironment();

  const onSubmit = async (data: IngressFormInput) => {
    await createIngressMutation.mutateAsync(
      {
        id: currentEnvironment.id,
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
        onSuccess: () => {
          queryClient.invalidateQueries({
            queryKey: trpc.environment.getEnvironmentBuilds.queryKey({
              id: currentEnvironment.id,
            }),
          });
        },
      },
    );
  };

  return (
    <DeployIngressForm
      deploymentEnvironment={deploymentEnvironment}
      resetOnSuccess={true}
      onSubmit={onSubmit}
    />
  );
}

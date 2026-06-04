import { useMutation, useQueryClient } from "@tanstack/react-query";
import DeployDatabaseForm from "~/components/organisms/forms/DeployDatabaseForm";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";

export default function Index() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();

  const { environment: currentEnvironment } = useEnvironment();

  const createDatabaseMutation = useMutation(
    trpc.deployment.deployDatabase.mutationOptions(),
  );

  const onSubmit = async (data: { serviceName: string }) => {
    if (!currentEnvironment) {
      return;
    }

    await createDatabaseMutation.mutateAsync(
      {
        id: currentEnvironment.id,
        serviceName: data.serviceName,
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

  return <DeployDatabaseForm resetOnSuccess onSubmit={onSubmit} />;
}

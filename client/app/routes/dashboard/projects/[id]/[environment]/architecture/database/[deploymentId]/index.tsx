import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import DeployDatabaseForm from "~/components/organisms/forms/DeployDatabaseForm";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";

export default function UpdateDatabaseDeployment() {
  const { slug, id, environment, deploymentId } = useParams<{
    slug: string;
    id: string;
    environment: string;
    deploymentId: string;
  }>();
  const navigate = useNavigate();

  const { environment: currentEnvironment } = useEnvironment();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const { data: environmentDeployments, isLoading } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: currentEnvironment?.id },
      { enabled: !!currentEnvironment },
    ),
  );

  const updateDatabaseMutation = useMutation(
    trpc.deployment.updateDatabase.mutationOptions(),
  );

  const onSubmit = async () => {
    if (!currentEnvironment) {
      return;
    }

    await updateDatabaseMutation.mutateAsync(
      {
        id: currentEnvironment.id,
        deploymentId: Number(deploymentId),
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
              `/${slug}/projects/${id}/${environment}/architecture/database/${result.deploymentId}`,
              { replace: true },
            );
          }
        },
      },
    );
  };

  const databaseDeployment = environmentDeployments?.databases.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployDatabaseForm
          key={deploymentId}
          onSubmit={onSubmit}
          defaultValues={{
            serviceName: databaseDeployment?.serviceName ?? "",
          }}
        />
      )}
    </>
  );
}

import React from "react";
import { useParams } from "react-router";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import DeployFromGitForm, {
  type DeployFromGitFormInput,
} from "~/components/organisms/forms/DeployFromGitForm";

export default function UpdateGitDeployment() {
  const { deploymentId } = useParams<{ deploymentId: string }>();

  const { environment: currentEnvironment, teamId } = useEnvironment();

  const trpc = useTRPC();
  const queryClient = useQueryClient();

  const { data: environmentDeployments, isLoading } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: currentEnvironment?.id },
      { enabled: !!currentEnvironment },
    ),
  );

  const updateGitRepoMutation = useMutation(
    trpc.deployment.updateDeployFromGitRepo.mutationOptions(),
  );

  const onSubmit = async (data: DeployFromGitFormInput) => {
    if (!data.port) {
      return;
    }

    await updateGitRepoMutation.mutateAsync(
      {
        id: currentEnvironment.id,
        deploymentId: Number(deploymentId),
        port: data.port,
        dockerfilePath: data.dockerfilePath,
        projectRepositoryPath: data.projectDirectoryPath,
        envs: data.envs,
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

  const gitDeployment = environmentDeployments?.gitDeployments.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployFromGitForm
          key={deploymentId}
          onSubmit={onSubmit}
          teamId={teamId}
          defaultValues={{
            serviceName: gitDeployment?.serviceName ?? "",
            url: gitDeployment?.gitUrl ?? "",
            dockerfilePath: gitDeployment?.dockerfilePath ?? "",
            projectDirectoryPath: gitDeployment?.projectRepositoryPath ?? "",
            port: gitDeployment ? Number(gitDeployment?.port) : null,
            envs: gitDeployment?.envVars ?? [],
          }}
        />
      )}
    </>
  );
}

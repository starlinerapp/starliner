import React from "react";
import DeployFromGitForm, {
  type DeployFromGitFormInput,
} from "~/components/organisms/forms/DeployFromGitForm";
import { useTRPC } from "~/utils/trpc/react";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useMutation } from "@tanstack/react-query";

export default function Git() {
  const trpc = useTRPC();

  const { environment: currentEnvironment } = useEnvironment();

  const createDeploymentMutation = useMutation(
    trpc.deployment.deployFromGitRepo.mutationOptions(),
  );

  const onSubmit = async (data: DeployFromGitFormInput) => {
    if (!data.port) {
      return;
    }

    await createDeploymentMutation.mutateAsync({
      environmentId: currentEnvironment.id,
      serviceName: data.serviceName,
      port: data.port,
      gitUrl: data.url,
      dockerfilePath: data.dockerfilePath,
      projectRepositoryPath: data.projectDirectoryPath,
      envs: data.envs,
    });
  };

  return (
    <>
      <DeployFromGitForm onSubmit={onSubmit} resetOnSuccess={true} />
    </>
  );
}

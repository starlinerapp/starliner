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
    await createDeploymentMutation.mutateAsync({
      environmentId: currentEnvironment.id,
      serviceName: data.serviceName,
      gitUrl: data.url,
      dockerfilePath: data.dockerfilePath,
      projectRepositoryPath: data.projectDirectoryPath,
    });
  };

  return (
    <>
      <DeployFromGitForm onSubmit={onSubmit} resetOnSuccess={true} />
    </>
  );
}

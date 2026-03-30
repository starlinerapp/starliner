import React from "react";
import DeployFromGitForm, {
  type DeployFromGitFormInput,
} from "~/components/organisms/forms/DeployFromGitForm";
import { useTRPC } from "~/utils/trpc/react";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import { useLoaderData } from "react-router";

export function loader() {
  return {
    githubAppName: process.env.GITHUB_APP_NAME,
  };
}

export default function Git() {
  const { githubAppName } = useLoaderData<typeof loader>();

  console.log(githubAppName);

  const trpc = useTRPC();
  const queryClient = useQueryClient();

  const { environment: currentEnvironment } = useEnvironment();
  const organization = useOrganizationContext();

  const { data: githubApp, isLoading } = useQuery(
    trpc.githubApp.getGithubApp.queryOptions({
      organizationId: organization.id,
    }),
  );

  const createDeploymentMutation = useMutation(
    trpc.deployment.deployFromGitRepo.mutationOptions(),
  );

  const onSubmit = async (data: DeployFromGitFormInput) => {
    if (!data.port) {
      return;
    }

    await createDeploymentMutation.mutateAsync(
      {
        environmentId: currentEnvironment.id,
        serviceName: data.serviceName,
        port: data.port,
        gitUrl: data.url,
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

  if (isLoading) {
    // TODO: Add Skeleton
    return null;
  }

  return (
    <>
      {githubApp ? (
        <DeployFromGitForm onSubmit={onSubmit} resetOnSuccess={true} />
      ) : (
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <p>Install GitHub App</p>
            <p className="text-mauve-11 truncate text-sm">
              Install the GitHub App to get started.
            </p>
          </div>
          <InstallGitHubApp githubAppName={githubAppName} />
        </div>
      )}
    </>
  );
}

import React from "react";
import DeployFromGitForm, {
  type DeployFromGitFormInput,
} from "~/components/organisms/forms/DeployFromGitForm";
import { useTRPC } from "~/utils/trpc/react";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import { useLoaderData, useLocation } from "react-router";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

export function loader() {
  return {
    githubAppName: process.env.GITHUB_APP_NAME,
  };
}

export default function Git() {
  const { githubAppName } = useLoaderData<typeof loader>();
  const location = useLocation();

  const trpc = useTRPC();
  const queryClient = useQueryClient();

  const { environment: currentEnvironment, teamId } = useEnvironment();
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
    return (
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <Skeleton className="h-6 w-44" />
          <Skeleton className="h-6 w-64" />
        </div>

        <div className="flex flex-col gap-2">
          <div className="flex flex-col gap-1">
            <Skeleton className="h-5 w-24" />
            <Skeleton className="h-9 w-full" />
          </div>

          <div className="flex flex-col gap-1">
            <Skeleton className="h-5 w-20" />
            <Skeleton className="h-9 w-full" />
          </div>

          <div className="flex items-end gap-2">
            <div className="flex w-full flex-col gap-1">
              <Skeleton className="h-5 w-28" />
              <Skeleton className="h-9 w-full" />
            </div>
            <div className="flex w-full flex-col gap-1">
              <Skeleton className="h-5 w-20" />
              <Skeleton className="h-9 w-full" />
            </div>
          </div>

          <div className="flex flex-col gap-1">
            <Skeleton className="h-5 w-10" />
            <Skeleton className="h-9 w-full" />
          </div>

          <div className="flex flex-col gap-1">
            <Skeleton className="h-5 w-36" />
            <Skeleton className="h-6 w-24" />
          </div>
        </div>

        <Skeleton className="h-8 w-28" />
      </div>
    );
  }

  return (
    <>
      {githubApp ? (
        <DeployFromGitForm
          onSubmit={onSubmit}
          resetOnSuccess={true}
          teamId={teamId}
        />
      ) : (
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <p>Install GitHub App</p>
            <p className="text-mauve-11 truncate text-sm">
              Install the GitHub App to get started.
            </p>
          </div>
          <InstallGitHubApp githubAppName={githubAppName} redirectTo={location.pathname} />
        </div>
      )}
    </>
  );
}

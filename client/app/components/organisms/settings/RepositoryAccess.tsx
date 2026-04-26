import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import React, { useState } from "react";
import { useLocation } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export function RepositoryAccess({
  teamId,
  githubAppName,
}: {
  teamId: number;
  githubAppName: string | undefined;
}) {
  const location = useLocation();
  const [showAssignDialog, setShowAssignDialog] = useState(false);
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();

  const { data: teamRepos, isLoading: isTeamReposLoading } = useQuery(
    trpc.team.getTeamRepositories.queryOptions({
      teamId,
    }),
  );

  const { data: githubApp, isLoading: isGithubAppLoading } = useQuery(
    trpc.githubApp.getGithubApp.queryOptions({
      organizationId: organization.id,
    }),
  );

  const { data: allRepos, isLoading: isAllReposLoading } = useQuery({
    ...trpc.github.getAllRepositories.queryOptions({
      organizationId: organization.id,
    }),
    enabled: organization.isOwner,
  });

  const assignMutation = useMutation(
    trpc.team.assignRepoToTeam.mutationOptions(),
  );

  const unassignMutation = useMutation(
    trpc.team.unassignRepoFromTeam.mutationOptions(),
  );

  function onAssignRepo(repoId: number, repoName: string) {
    assignMutation.mutate(
      {
        teamId,
        githubRepoId: repoId,
        repoName,
      },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamRepositories.queryKey(),
          });
        },
      },
    );
  }

  function onUnassignRepo(repoId: number) {
    unassignMutation.mutate(
      {
        teamId,
        githubRepoId: repoId,
      },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamRepositories.queryKey(),
          });
        },
      },
    );
  }

  const assignedRepoIds = new Set(teamRepos?.map((r) => r.githubRepoId) ?? []);
  const allReposSorted =
    allRepos
      ?.slice()
      .sort((a, b) =>
        a.name.localeCompare(b.name, undefined, { sensitivity: "base" }),
      ) ?? [];

  return (
    <div className="w-full">
      <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
          <p>Repository Access</p>
          {githubApp && (
            <Dialog open={showAssignDialog} onOpenChange={setShowAssignDialog}>
              <DialogTrigger asChild>
                <Button intent="secondary" className="h-7 w-36 text-xs">
                  Manage Repositories
                </Button>
              </DialogTrigger>
              <DialogContent>
                <div className="flex flex-col gap-4">
                  <div className="flex flex-col gap-2">
                    <h1>Manage repository access</h1>
                    <p className="text-mauve-11 text-sm">
                      Add or remove repositories to control what this team can
                      see.
                    </p>
                  </div>
                  {isAllReposLoading ? (
                    <div className="flex flex-col gap-2">
                      <Skeleton className="h-8 w-full" />
                      <Skeleton className="h-8 w-full" />
                      <Skeleton className="h-8 w-full" />
                    </div>
                  ) : allReposSorted.length === 0 ? (
                    <div className="text-mauve-11 text-sm">
                      No GitHub repositories are available.
                    </div>
                  ) : (
                    <div className="border-mauve-6 divide-mauve-6 max-h-[60vh] divide-y overflow-y-auto rounded-md border">
                      {allReposSorted.map((repo) => {
                        const isAssigned = assignedRepoIds.has(repo.id);
                        return (
                          <div
                            key={repo.id}
                            className="flex min-w-0 items-center justify-between gap-3 px-4 py-2"
                          >
                            <div className="flex min-w-0 flex-1 flex-col">
                              <span className="text-mauve-12 truncate text-sm font-medium">
                                {repo.owner}/{repo.name}
                              </span>
                              {repo.description && (
                                <span
                                  className="text-mauve-11 truncate text-xs"
                                  title={repo.description}
                                >
                                  {repo.description}
                                </span>
                              )}
                            </div>
                            {isAssigned ? (
                              <Button
                                className="h-7 w-24 text-xs"
                                intent="secondary"
                                onClick={() => onUnassignRepo(repo.id)}
                              >
                                Unassign
                              </Button>
                            ) : (
                              <Button
                                className="h-7 w-24 text-xs"
                                intent="primary"
                                onClick={() => {
                                  onAssignRepo(repo.id, repo.full_name);
                                }}
                              >
                                Assign
                              </Button>
                            )}
                          </div>
                        );
                      })}
                    </div>
                  )}
                </div>
              </DialogContent>
            </Dialog>
          )}
        </div>
        {isGithubAppLoading ? (
          <div className="flex flex-col gap-2 px-4 py-3">
            <Skeleton className="h-5 w-48" />
          </div>
        ) : !githubApp ? (
          <div className="flex flex-col gap-4 px-4 py-3">
            <div className="flex flex-col gap-1">
              <p className="text-mauve-11 text-sm">
                Install the GitHub App to assign repositories to this team.
              </p>
            </div>
            <InstallGitHubApp
              githubAppName={githubAppName}
              redirectTo={location.pathname}
            />
          </div>
        ) : isTeamReposLoading ? (
          <div className="flex flex-col gap-2 px-4 py-3">
            <Skeleton className="h-5 w-48" />
          </div>
        ) : teamRepos?.length === 0 ? (
          <div className="text-mauve-11 px-4 py-3 text-sm">
            No repositories assigned. Team members cannot see any repositories
            until you assign them.
          </div>
        ) : (
          teamRepos?.map((repo) => (
            <div
              key={repo.githubRepoId}
              className="border-mauve-6 text-mauve-12 min-w-0 border-b px-4 py-3 text-sm last:border-b-0"
            >
              <span className="block truncate" title={repo.repoName}>
                {repo.repoName}
              </span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import { Cross } from "~/components/atoms/icons";
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
  const unassignedRepos =
    allRepos?.filter((r) => !assignedRepoIds.has(r.id)) ?? [];

  return (
    <div className="w-full">
      <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
          <p>Repository Access</p>
          {githubApp && (
            <Dialog open={showAssignDialog} onOpenChange={setShowAssignDialog}>
              <DialogTrigger asChild>
                <Button intent="secondary" className="h-7 w-28 text-xs">
                  Assign Repo
                </Button>
              </DialogTrigger>
              <DialogContent>
                <div className="flex flex-col gap-4">
                  <div className="flex flex-col gap-2">
                    <h1>Assign Repository</h1>
                    <p className="text-mauve-11 text-sm">
                      Select a repository to make visible to this team&apos;s
                      members.
                    </p>
                  </div>
                  {isAllReposLoading ? (
                    <div className="flex flex-col gap-2">
                      <Skeleton className="h-8 w-full" />
                      <Skeleton className="h-8 w-full" />
                      <Skeleton className="h-8 w-full" />
                    </div>
                  ) : unassignedRepos.length === 0 ? (
                    <div className="text-mauve-11 text-sm">
                      All repositories are already assigned to this team.
                    </div>
                  ) : (
                    <div className="border-mauve-6 divide-mauve-6 max-h-[60vh] divide-y overflow-y-auto rounded-md border">
                      {unassignedRepos.map((repo) => (
                        <div
                          key={repo.id}
                          className="flex items-center justify-between px-2 py-2"
                        >
                          <div className="flex flex-col">
                            <span className="text-mauve-12 text-sm font-medium">
                              {repo.name}
                            </span>
                            {repo.description && (
                              <span className="text-mauve-11 text-xs">
                                {repo.description}
                              </span>
                            )}
                          </div>
                          <Button
                            className="h-7 w-20 text-xs"
                            intent="secondary"
                            disabled={assignMutation.isPending}
                            onClick={() => {
                              onAssignRepo(repo.id, repo.full_name);
                            }}
                          >
                            Assign
                          </Button>
                        </div>
                      ))}
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
              className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
            >
              <span className="text-sm">{repo.repoName}</span>
              <button
                className="text-mauve-11 hover:bg-gray-3 cursor-pointer rounded-md p-1"
                disabled={unassignMutation.isPending}
                onClick={() => onUnassignRepo(repo.githubRepoId)}
                title="Revoke repository access"
              >
                <Cross width={20} height={20} />
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

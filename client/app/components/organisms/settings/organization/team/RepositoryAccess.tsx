import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { useLocation } from "react-router";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import Button from "~/components/atoms/button/Button";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

export function RepositoryAccess({
  teamId,
  githubAppName,
}: {
  teamId: number;
  githubAppName: string | undefined;
}) {
  const location = useLocation();
  const [showAssignDialog, setShowAssignDialog] = useState(false);
  const [pendingAssignedRepoIds, setPendingAssignedRepoIds] = useState<
    Set<number>
  >(new Set());

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
    enabled: organization.isOwner && !!githubApp,
  });

  const setTeamRepositoriesMutation = useMutation(
    trpc.team.setTeamRepositories.mutationOptions(),
  );

  const allReposSorted =
    allRepos?.slice().sort((a, b) =>
      a.name.localeCompare(b.name, undefined, {
        sensitivity: "base",
      }),
    ) ?? [];

  function getAssignedRepoIds() {
    return new Set(teamRepos?.map((r) => r.githubRepoId) ?? []);
  }

  useEffect(() => {
    if (showAssignDialog) {
      setPendingAssignedRepoIds(
        new Set(teamRepos?.map((r) => r.githubRepoId) ?? []),
      );
    }
  }, [showAssignDialog, teamRepos]);

  function toggleRepo(repoId: number, checked: boolean) {
    setPendingAssignedRepoIds((prev) => {
      const next = new Set(prev);

      if (checked) {
        next.add(repoId);
      } else {
        next.delete(repoId);
      }

      return next;
    });
  }

  function onApply() {
    const repositories = allReposSorted
      .filter((repo) => pendingAssignedRepoIds.has(repo.id))
      .map((repo) => ({
        githubRepoId: repo.id,
        repoName: `${repo.owner}/${repo.name}`,
      }));

    setTeamRepositoriesMutation.mutate(
      {
        teamId,
        repositories,
      },
      {
        onSuccess: async () => {
          await Promise.all([
            queryClient.invalidateQueries({
              queryKey: trpc.team.getTeamRepositories.queryKey(),
            }),
            queryClient.invalidateQueries({
              queryKey: trpc.github.getRepositories.queryKey(),
            }),
          ]);

          setShowAssignDialog(false);
        },
      },
    );
  }

  function onCancel() {
    setPendingAssignedRepoIds(getAssignedRepoIds());
    setShowAssignDialog(false);
  }

  return (
    <div className="w-full">
      <div className="rounded-md border border-mauve-6 text-sm shadow-xs">
        <div className="flex h-14 items-center justify-between rounded-t-md border-mauve-6 border-b bg-gray-2 px-4 font-bold text-mauve-12 text-xs uppercase">
          <p>Repositories</p>
          {organization.isOwner && (
            <Dialog
              open={showAssignDialog}
              onOpenChange={(open) => {
                setShowAssignDialog(open);
                if (!open) {
                  setPendingAssignedRepoIds(getAssignedRepoIds());
                }
              }}
            >
              <DialogTrigger asChild>
                <Button intent="secondary" className="w-36 text-xs">
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
                      <Skeleton className="h-12 w-full" />
                      <Skeleton className="h-12 w-full" />
                      <Skeleton className="h-12 w-full" />
                    </div>
                  ) : allReposSorted.length === 0 ? (
                    <WarningBanner text="Install the GitHub App to assign repositories to this team." />
                  ) : (
                    <div className="flex max-h-[60vh] flex-col overflow-y-auto rounded-md border border-mauve-6 bg-mauve-2">
                      {allReposSorted.map((repo) => (
                        <label
                          key={repo.id}
                          className="flex min-w-0 cursor-pointer items-center gap-3 p-3"
                        >
                          <input
                            type="checkbox"
                            checked={pendingAssignedRepoIds.has(repo.id)}
                            onChange={(event) => {
                              toggleRepo(repo.id, event.target.checked);
                            }}
                            className="h-4.5 w-4.5 shrink-0 rounded border-mauve-6"
                          />
                          <div className="flex min-w-0 flex-1 flex-col gap-1">
                            <p className="truncate font-medium text-mauve-12 text-sm">
                              {repo.owner}/{repo.name}
                            </p>
                            {repo.description && (
                              <p
                                className="truncate text-mauve-11 text-xs"
                                title={repo.description}
                              >
                                {repo.description}
                              </p>
                            )}
                          </div>
                        </label>
                      ))}
                    </div>
                  )}
                  <div className="flex w-full justify-end gap-2">
                    <Button
                      intent="secondary"
                      className="w-24"
                      onClick={onCancel}
                      disabled={setTeamRepositoriesMutation.isPending}
                    >
                      Cancel
                    </Button>
                    <Button
                      className="w-24"
                      onClick={onApply}
                      disabled={
                        setTeamRepositoriesMutation.isPending ||
                        allReposSorted.length === 0
                      }
                    >
                      Apply
                    </Button>
                  </div>
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
            {organization.isOwner && (
              <InstallGitHubApp
                githubAppName={githubAppName}
                redirectTo={location.pathname}
              />
            )}
          </div>
        ) : isTeamReposLoading ? (
          <div className="flex flex-col gap-2 px-4 py-3">
            <Skeleton className="h-5 w-48" />
          </div>
        ) : teamRepos?.length === 0 ? (
          <p className="px-4 py-3 text-mauve-11 text-sm">
            No repositories assigned.
          </p>
        ) : (
          teamRepos?.map((repo) => (
            <div
              key={repo.githubRepoId}
              className="min-w-0 border-mauve-6 border-b px-4 py-3 text-mauve-12 text-sm last:border-b-0"
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

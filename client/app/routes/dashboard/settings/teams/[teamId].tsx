import React, { useState } from "react";
import { useParams } from "react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import Button from "~/components/atoms/button/Button";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import { Trash } from "~/components/atoms/icons";

export default function TeamDetail() {
  const { teamId } = useParams();
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();
  const [showAssignDialog, setShowAssignDialog] = useState(false);

  const { data: members, isLoading: isMembersLoading } = useQuery(
    trpc.team.getTeamMembers.queryOptions({
      teamId: Number(teamId),
    }),
  );

  const { data: teamRepos, isLoading: isTeamReposLoading } = useQuery(
    trpc.team.getTeamRepositories.queryOptions({
      organizationId: organization.id,
      teamId: Number(teamId),
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
        organizationId: organization.id,
        teamId: Number(teamId),
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
        organizationId: organization.id,
        teamId: Number(teamId),
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

  const assignedRepoIds = new Set(
    teamRepos?.map((r) => r.github_repo_id) ?? [],
  );
  const unassignedRepos =
    allRepos?.filter((r) => !assignedRepoIds.has(r.id)) ?? [];

  return (
    <div className="flex flex-col gap-8 px-8 py-4">
      <h1 className="pt-1 text-xl font-bold">Team Settings</h1>

      <div className="w-full xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
            Members
          </div>
          {isMembersLoading ? (
            <div className="flex flex-col gap-2 px-4 py-3">
              <Skeleton className="h-5 w-48" />
              <Skeleton className="h-5 w-36" />
              <Skeleton className="h-5 w-52" />
            </div>
          ) : members?.length === 0 ? (
            <div className="text-mauve-11 px-4 py-3 text-sm">
              No members yet.
            </div>
          ) : (
            members?.map((member) => (
              <div
                key={member.user_id}
                className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
              >
                <span>{member.name}</span>
                <span className="text-mauve-11 text-xs">{member.email}</span>
              </div>
            ))
          )}
        </div>
      </div>

      {organization.isOwner && (
        <div className="w-full xl:w-3/5">
          <div className="border-mauve-6 rounded-md border-1 text-sm">
            <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
              <p>Repository Access</p>
              <Dialog
                open={showAssignDialog}
                onOpenChange={setShowAssignDialog}
              >
                <DialogTrigger asChild>
                  <Button className="h-7 w-28 text-xs">Assign Repo</Button>
                </DialogTrigger>
                <DialogContent>
                  <h2 className="text-mauve-12 mb-2 text-lg font-bold">
                    Assign Repository
                  </h2>
                  <p className="text-mauve-11 mb-4 text-sm">
                    Select a repository to make visible to this team&apos;s
                    members.
                  </p>
                  {isAllReposLoading ? (
                    <div className="flex flex-col gap-2 py-3">
                      <Skeleton className="h-8 w-full" />
                      <Skeleton className="h-8 w-full" />
                      <Skeleton className="h-8 w-full" />
                    </div>
                  ) : unassignedRepos.length === 0 ? (
                    <div className="text-mauve-11 py-3 text-sm">
                      All repositories are already assigned to this team.
                    </div>
                  ) : (
                    <div className="max-h-80 overflow-y-auto">
                      {unassignedRepos.map((repo) => (
                        <div
                          key={repo.id}
                          className="border-mauve-6 flex items-center justify-between border-b px-2 py-2 last:border-b-0"
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
                </DialogContent>
              </Dialog>
            </div>
            {isTeamReposLoading ? (
              <div className="flex flex-col gap-2 px-4 py-3">
                <Skeleton className="h-5 w-48" />
                <Skeleton className="h-5 w-36" />
              </div>
            ) : teamRepos?.length === 0 ? (
              <div className="text-mauve-11 px-4 py-3 text-sm">
                No repositories assigned. Team members cannot see any
                repositories until you assign them.
              </div>
            ) : (
              teamRepos?.map((repo) => (
                <div
                  key={repo.github_repo_id}
                  className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
                >
                  <span className="font-mono text-xs">{repo.repo_name}</span>
                  <button
                    className="text-red-11 hover:text-red-9 cursor-pointer"
                    disabled={unassignMutation.isPending}
                    onClick={() => onUnassignRepo(repo.github_repo_id)}
                    title="Remove repository access"
                  >
                    <Trash width={16} height={16} />
                  </button>
                </div>
              ))
            )}
          </div>
          <p className="text-mauve-11 mt-2 text-xs">
            Only repositories explicitly assigned to a team will be visible to
            its members. By default, no repositories are accessible.
          </p>
        </div>
      )}
    </div>
  );
}

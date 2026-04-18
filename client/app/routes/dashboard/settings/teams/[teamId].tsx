import React, { useMemo, useState } from "react";
import {
  useNavigate,
  useParams,
  useLoaderData,
  useLocation,
} from "react-router";
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
import { ChevronDown, MagnifyingGlass, Cross } from "~/components/atoms/icons";
import * as Popover from "@radix-ui/react-popover";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";

export function loader() {
  return {
    githubAppName: process.env.GITHUB_APP_NAME,
  };
}

export default function TeamDetail() {
  const { githubAppName } = useLoaderData<typeof loader>();
  const { teamId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();
  const [showAssignDialog, setShowAssignDialog] = useState(false);
  const [search, setSearch] = useState("");
  const [addMemberOpen, setAddMemberOpen] = useState(false);

  const { data: user } = useQuery(trpc.user.getUser.queryOptions());

  const { data: teams } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );
  const team = teams?.find((t) => t.id === Number(teamId));

  const { data: members, isLoading } = useQuery(
    trpc.team.getTeamMembers.queryOptions({
      teamId: Number(teamId),
    }),
  );

  const { data: orgMembers } = useQuery({
    ...trpc.organization.getOrganizationMembers.queryOptions({
      id: organization.id,
    }),
    enabled: organization.isOwner,
  });

  const addMemberMutation = useMutation(
    trpc.team.addTeamMember.mutationOptions(),
  );

  const removeMemberMutation = useMutation(
    trpc.team.removeTeamMember.mutationOptions(),
  );

  const filteredOrgMembers = useMemo(() => {
    if (!orgMembers || !members) return [];
    const memberIds = new Set(members.map((m) => m.user_id));
    return orgMembers.filter(
      (m) =>
        !memberIds.has(m.user_id) &&
        (m.name.toLowerCase().includes(search.toLowerCase()) ||
          m.email.toLowerCase().includes(search.toLowerCase())),
    );
  }, [orgMembers, members, search]);

  function handleAddMember(userId: number) {
    addMemberMutation.mutate(
      { teamId: Number(teamId), userId },
      {
        onSuccess: async () => {
          setSearch("");
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamMembers.queryKey({
              teamId: Number(teamId),
            }),
          });
        },
      },
    );
  }

  function handleRemoveMember(userId: number) {
    removeMemberMutation.mutate(
      { teamId: Number(teamId), userId },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamMembers.queryKey({
              teamId: Number(teamId),
            }),
          });
        },
      },
    );
  }

  const { data: teamRepos, isLoading: isTeamReposLoading } = useQuery(
    trpc.team.getTeamRepositories.queryOptions({
      organizationId: organization.id,
      teamId: Number(teamId),
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
      <div className="flex flex-col gap-1">
        <nav className="text-mauve-11 flex items-center gap-1 pb-4 text-xs">
          <button
            onClick={() => navigate("../", { relative: "path" })}
            className="hover:text-mauve-12 cursor-pointer hover:underline"
          >
            Teams
          </button>
          <span>&gt;</span>
          {isLoading ? (
            <Skeleton className="h-4 w-16" />
          ) : (
            <span className="text-mauve-12">#{team?.slug}</span>
          )}
        </nav>
        {isLoading ? (
          <Skeleton className="h-6 w-32" />
        ) : (
          <h1 className="text-xl font-bold">#{team?.slug}</h1>
        )}
      </div>
      <h1 className="pt-1 text-xl font-bold">Team</h1>
      <div className="flex flex-col gap-4">
        <div className="w-full xl:w-3/5">
          <div className="border-mauve-6 rounded-md border text-sm">
            <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
              <p>Members</p>
              {organization.isOwner && (
                <Popover.Root
                  open={addMemberOpen}
                  onOpenChange={setAddMemberOpen}
                >
                  <Popover.Trigger asChild>
                    <Button intent="secondary" className="h-7 w-28 text-xs">
                      Add Member
                      <ChevronDown
                        className={`h-3 w-3 ${addMemberOpen ? "rotate-180" : ""}`}
                      />
                    </Button>
                  </Popover.Trigger>
                  <Popover.Portal>
                    <Popover.Content className="border-mauve-6 w-70 space-y-2 rounded-md border bg-white p-2 shadow-md">
                      <div className="relative">
                        <MagnifyingGlass className="text-mauve-9 absolute top-1/2 left-2 h-4 w-4 -translate-y-1/2" />
                        <input
                          className="border-mauve-6 w-full rounded-md border p-2 pl-8 text-xs"
                          type="text"
                          placeholder="Search Members"
                          autoFocus
                          value={search}
                          onChange={(e) => setSearch(e.target.value)}
                        />
                      </div>
                      <div className="divide-gray-4 flex max-h-60 flex-col divide-y overflow-y-auto">
                        {filteredOrgMembers.length > 0 ? (
                          filteredOrgMembers.map((m) => (
                            <button
                              key={m.user_id}
                              onClick={() => handleAddMember(m.user_id)}
                              className="hover:bg-gray-3 flex w-full items-center justify-between rounded px-2 py-2 text-left transition-colors"
                            >
                              <span className="text-xs">{m.name}</span>
                              <span className="text-mauve-11 text-xs">
                                {m.email}
                              </span>
                            </button>
                          ))
                        ) : (
                          <div className="text-mauve-11 p-4 text-center text-xs">
                            No members found.
                          </div>
                        )}
                      </div>
                    </Popover.Content>
                  </Popover.Portal>
                </Popover.Root>
              )}
            </div>
            {isLoading ? (
              <>
                {Array.from({ length: 5 }).map((_, i) => (
                  <div
                    key={i}
                    className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-2 text-sm last:border-b-0"
                  >
                    <div className="flex flex-col gap-1">
                      <Skeleton className="h-5 w-24" />
                      <Skeleton className="h-5 w-36" />
                    </div>
                  </div>
                ))}
              </>
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
                  <div className="flex flex-col">
                    <span>{member.name}</span>
                    <span className="text-mauve-11">{member.email}</span>
                  </div>
                  {organization.isOwner &&
                    member.user_id !== Number(user?.user_id) && (
                      <Button
                        className="h-7 w-20 text-xs"
                        intent="secondary"
                        disabled={removeMemberMutation.isPending}
                        onClick={() => handleRemoveMember(member.user_id)}
                      >
                        Remove
                      </Button>
                    )}
                </div>
              ))
            )}
          </div>
        </div>

        {organization.isOwner && (
          <div className="w-full xl:w-3/5">
            <div className="border-mauve-6 rounded-md border text-sm">
              <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
                <p>Repository Access</p>
                {githubApp && (
                  <Dialog
                    open={showAssignDialog}
                    onOpenChange={setShowAssignDialog}
                  >
                    <DialogTrigger asChild>
                      <Button intent="secondary" className="h-7 w-28 text-xs">
                        Assign Repo
                      </Button>
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
                              className="border-mauve-6 group flex items-center justify-between border-b px-2 py-2 last:border-b-0"
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
                                className="h-7 w-20 text-xs opacity-0 transition-opacity group-hover:opacity-100"
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
                )}
              </div>
              {isGithubAppLoading ? (
                <div className="flex flex-col gap-2 px-4 py-3">
                  <Skeleton className="h-5 w-48" />
                  <Skeleton className="h-5 w-36" />
                </div>
              ) : !githubApp ? (
                <div className="flex flex-col gap-4 px-4 py-3">
                  <div className="flex flex-col gap-1">
                    <p className="text-mauve-11 text-xs">
                      Install the GitHub App to assign repositories to this
                      team.
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
                      className="text-mauve-11 hover:text-mauve-12 cursor-pointer"
                      disabled={unassignMutation.isPending}
                      onClick={() => onUnassignRepo(repo.github_repo_id)}
                      title="Revoke repository access"
                    >
                      <Cross width={16} height={16} />
                    </button>
                  </div>
                ))
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

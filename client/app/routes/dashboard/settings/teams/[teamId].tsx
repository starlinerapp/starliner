import React, { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import Button from "~/components/atoms/button/Button";
import { ChevronDown, MagnifyingGlass } from "~/components/atoms/icons";
import * as Popover from "@radix-ui/react-popover";

export default function TeamDetail() {
  const { teamId } = useParams();
  const navigate = useNavigate();
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();
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
      <div className="w-full xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
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
    </div>
  );
}

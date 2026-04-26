import * as Popover from "@radix-ui/react-popover";
import Button from "~/components/atoms/button/Button";
import { ChevronDown, MagnifyingGlass } from "~/components/atoms/icons";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import React, { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { cn } from "~/utils/cn";

export default function TeamMembers({ teamId }: { teamId: number }) {
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();

  const [search, setSearch] = useState("");
  const [addMemberOpen, setAddMemberOpen] = useState(false);

  const { data: user } = useQuery(trpc.user.getUser.queryOptions());

  const { data: members, isLoading } = useQuery(
    trpc.team.getTeamMembers.queryOptions({
      teamId,
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
      { teamId, userId },
      {
        onSuccess: async () => {
          setSearch("");
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamMembers.queryKey({
              teamId,
            }),
          });
        },
      },
    );
  }

  function handleRemoveMember(userId: number) {
    removeMemberMutation.mutate(
      { teamId, userId },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamMembers.queryKey({
              teamId,
            }),
          });
        },
      },
    );
  }

  return (
    <div className="border-mauve-6 rounded-md border text-sm">
      <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
        <p>Members</p>
        {organization.isOwner && (
          <Popover.Root open={addMemberOpen} onOpenChange={setAddMemberOpen}>
            <Popover.Trigger asChild>
              <Button intent="secondary" className="h-7 w-28 text-xs">
                Add Member
                <ChevronDown
                  className={cn("h-3 w-3", addMemberOpen && "rotate-180")}
                />
              </Button>
            </Popover.Trigger>
            <Popover.Portal>
              <Popover.Content className="border-mauve-6 mt-1 w-70 space-y-2 rounded-md border bg-white p-2 shadow-md">
                <div className="relative">
                  <MagnifyingGlass className="text-mauve-11 absolute top-1/2 left-2 h-4 w-4 -translate-y-1/2" />
                  <input
                    className="border-mauve-6 placeholder:text-mauve-11 w-full rounded-md border p-2 pl-7 text-xs"
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
                        <span className="text-mauve-11 text-xs">{m.email}</span>
                      </button>
                    ))
                  ) : (
                    <div className="text-mauve-11 p-4 text-center text-xs">
                      No members found
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
          {Array.from({ length: 3 }).map((_, i) => (
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
        <div className="text-mauve-11 px-4 py-3 text-sm">No members yet.</div>
      ) : (
        members?.map((member) => (
          <div
            key={member.user_id}
            className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm shadow-xs last:border-b-0"
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
  );
}

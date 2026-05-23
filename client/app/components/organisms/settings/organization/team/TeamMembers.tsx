import * as Popover from "@radix-ui/react-popover";
import Button from "~/components/atoms/button/Button";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { ChevronDown, MagnifyingGlass } from "~/components/atoms/icons";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import React, { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { cn } from "~/utils/cn";
import { AvatarIcon } from "~/components/atoms/avatar/Avatar";
import AddMemberDialog from "~/components/organisms/dialog/AddMemberDialog";

interface MemberToRemove {
  userId: number;
  name: string;
}

export default function TeamMembers({ teamId }: { teamId: number }) {
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();

  const [search, setSearch] = useState("");
  const [addMemberOpen, setAddMemberOpen] = useState(false);
  const [showInviteMemberDialog, setShowInviteMemberDialog] = useState(false);
  const [memberToRemove, setMemberToRemove] = useState<MemberToRemove | null>(
    null,
  );

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
        !memberIds.has(m.id) &&
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

  function openRemoveMemberDialog(member: MemberToRemove) {
    removeMemberMutation.reset();
    setMemberToRemove(member);
  }

  function confirmRemoveMember() {
    if (memberToRemove == null) return;

    removeMemberMutation.mutate(
      { teamId, userId: memberToRemove.userId },
      {
        onSuccess: async () => {
          setMemberToRemove(null);
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
    <>
      <div className="border-mauve-6 w-full rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center justify-between border-b px-4 text-xs font-bold uppercase">
          <p>Members</p>
          {organization.isOwner && (
            <Popover.Root open={addMemberOpen} onOpenChange={setAddMemberOpen}>
              <Popover.Trigger asChild>
                <Button intent="secondary" className="w-30 text-xs">
                  Add Member
                  <ChevronDown
                    className={cn("h-3 w-3", addMemberOpen && "rotate-180")}
                  />
                </Button>
              </Popover.Trigger>
              <Popover.Portal>
                <Popover.Content className="border-mauve-6 mx-2 my-1 w-70 rounded-md border bg-white shadow-md">
                  <div className="space-y-2 p-2">
                    <div className="relative">
                      <MagnifyingGlass className="text-mauve-11 absolute top-1/2 left-2 h-4 w-4 -translate-y-1/2" />
                      <input
                        className="border-mauve-6 placeholder:text-mauve-11 w-full rounded-md border p-2 pl-7 text-xs shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
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
                            key={m.id}
                            onClick={() => handleAddMember(m.id)}
                            className="hover:bg-gray-3 flex w-full cursor-pointer items-center gap-3 rounded p-2 text-left transition-colors"
                          >
                            <AvatarIcon
                              name={m.name}
                              profilePicture={m.avatarUrl}
                            />
                            <span>
                              <p className="text-bold text-mauve-12 text-xs">
                                {m.name}
                              </p>
                              <p className="text-mauve-11 text-xs">{m.email}</p>
                            </span>
                          </button>
                        ))
                      ) : (
                        <div className="text-mauve-11 flex h-12 items-center justify-center text-xs">
                          <p>No members found</p>
                        </div>
                      )}
                    </div>
                  </div>
                  <div className="border-mauve-6 bg-mauve-2 border-t p-2 text-xs">
                    <Button
                      type="button"
                      intent="primary"
                      className="text-xs"
                      onClick={() => {
                        setAddMemberOpen(false);
                        setShowInviteMemberDialog(true);
                      }}
                    >
                      Invite Member
                    </Button>
                  </div>
                </Popover.Content>
              </Popover.Portal>
            </Popover.Root>
          )}
        </div>
        {isLoading ? (
          <>
            {Array.from({ length: 1 }).map((_, i) => (
              <div
                key={i}
                className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
              >
                <div className="flex items-center gap-3">
                  <Skeleton className="h-8 w-8 rounded-full" />
                  <div className="flex flex-col gap-1">
                    <Skeleton className="h-4.5 w-28" />
                    <Skeleton className="h-4.5 w-44" />
                  </div>
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
              className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
            >
              <div className="flex items-center gap-3">
                <AvatarIcon
                  name={member.name}
                  profilePicture={member.avatarUrl}
                />
                <div className="flex flex-col">
                  <span>{member.name}</span>
                  <span className="text-mauve-11">{member.email}</span>
                </div>
              </div>
              {organization.isOwner &&
                member.user_id !== Number(user?.user_id) && (
                  <Button
                    className="text-mauve-12 w-20 text-xs"
                    intent="secondary"
                    disabled={removeMemberMutation.isPending}
                    onClick={() =>
                      openRemoveMemberDialog({
                        userId: member.user_id,
                        name: member.name,
                      })
                    }
                  >
                    Remove
                  </Button>
                )}
            </div>
          ))
        )}
      </div>

      {organization.isOwner && (
        <AddMemberDialog
          organizationId={organization.id}
          teamId={teamId}
          open={showInviteMemberDialog}
          onOpenChange={setShowInviteMemberDialog}
        />
      )}

      <Dialog
        open={memberToRemove != null}
        onOpenChange={(open) => {
          if (!open) {
            setMemberToRemove(null);
            removeMemberMutation.reset();
          }
        }}
      >
        <DialogContent>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <h1>Remove Member</h1>
              <p className="text-mauve-11 text-sm">
                Are you sure you want to remove{" "}
                <span className="text-mauve-12 font-medium">
                  {memberToRemove?.name}
                </span>{" "}
                from this team?
              </p>
            </div>
            {removeMemberMutation.isError && (
              <ErrorBanner text={removeMemberMutation.error.message} />
            )}
            <div className="flex justify-end gap-2">
              <Button
                type="button"
                intent="secondary"
                className="w-24 cursor-pointer"
                disabled={removeMemberMutation.isPending}
                onClick={() => setMemberToRemove(null)}
              >
                Cancel
              </Button>
              <Button
                className="w-24 cursor-pointer"
                intent="primary"
                disabled={removeMemberMutation.isPending}
                onClick={confirmRemoveMember}
              >
                Remove
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

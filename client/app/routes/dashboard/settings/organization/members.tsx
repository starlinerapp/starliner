import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { AvatarIcon } from "~/components/atoms/avatar/Avatar";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import Button from "~/components/atoms/button/Button";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import AddMemberDialog from "~/components/organisms/dialog/AddMemberDialog";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

interface MemberToRemove {
  userId: number;
  name: string;
}

export default function Members() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const organization = useOrganizationContext();
  const [showAddMemberDialog, setShowAddMemberDialog] = useState(false);
  const [memberToRemove, setMemberToRemove] = useState<MemberToRemove | null>(
    null,
  );

  const { data: user } = useQuery(trpc.user.getUser.queryOptions());

  const { data: members, isLoading } = useQuery(
    trpc.organization.getOrganizationMembers.queryOptions({
      id: organization.id,
    }),
  );

  const removeMemberMutation = useMutation(
    trpc.organization.removeOrganizationMember.mutationOptions(),
  );

  function openRemoveMemberDialog(member: MemberToRemove) {
    removeMemberMutation.reset();
    setMemberToRemove(member);
  }

  function confirmRemoveMember() {
    if (memberToRemove == null) return;

    removeMemberMutation.mutate(
      {
        organizationId: organization.id,
        userId: memberToRemove.userId,
      },
      {
        onSuccess: async () => {
          setMemberToRemove(null);
          await queryClient.invalidateQueries({
            queryKey: trpc.organization.getOrganizationMembers.queryKey({
              id: organization.id,
            }),
          });
        },
      },
    );
  }

  return (
    <>
      <Breadcrumbs
        crumbs={[
          { label: "Settings" },
          { label: "Organization" },
          { label: "Members" },
        ]}
      />
      <div className="flex flex-col px-4 py-4">
        <div className="rounded-md border border-mauve-6 bg-gray-2 text-sm shadow-xs">
          <div className="flex h-14 items-center justify-between rounded-t-md px-4 font-bold text-mauve-12 text-xs uppercase">
            Members
            {organization.isOwner && (
              <Button
                className="w-28 text-xs"
                intent="secondary"
                onClick={() => setShowAddMemberDialog(true)}
              >
                Invite Member
              </Button>
            )}
          </div>
          <div className="mx-1 mb-1 divide-y divide-mauve-6 overflow-hidden rounded-md border border-mauve-6 bg-white shadow-xs">
            {isLoading ? (
              Array.from({ length: 2 }).map((_, i) => (
                <div
                  key={i}
                  className="flex h-14 items-center justify-between gap-2 px-4"
                >
                  <div className="flex items-center gap-3">
                    <Skeleton className="h-8 w-8 rounded-full" />
                    <div className="flex flex-col gap-1">
                      <Skeleton className="h-4 w-24" />
                      <Skeleton className="h-4 w-36" />
                    </div>
                  </div>
                  <Skeleton className="h-4 w-16" />
                </div>
              ))
            ) : members?.length === 0 ? (
              <div className="flex h-14 items-center px-4 text-mauve-11 text-sm">
                No members yet.
              </div>
            ) : (
              members?.map((member) => (
                <div
                  key={member.id}
                  className="flex h-14 items-center justify-between gap-2 px-4"
                >
                  <div className="flex items-center gap-3">
                    <AvatarIcon
                      name={member.name}
                      profilePicture={member.avatarUrl}
                    />
                    <div className="flex flex-col">
                      <span className="font-medium text-mauve-12">
                        {member.name}
                      </span>
                      <span className="text-mauve-11 text-sm">
                        {member.email}
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-mauve-11">
                      {member.is_owner ? "Owner" : "Member"}
                    </span>
                    {organization.isOwner &&
                      !member.is_owner &&
                      member.id !== Number(user?.user_id) && (
                        <Button
                          className="w-20 text-mauve-12 text-xs"
                          intent="secondary"
                          disabled={removeMemberMutation.isPending}
                          onClick={() =>
                            openRemoveMemberDialog({
                              userId: member.id,
                              name: member.name,
                            })
                          }
                        >
                          Remove
                        </Button>
                      )}
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {organization.isOwner && (
        <AddMemberDialog
          organizationId={organization.id}
          open={showAddMemberDialog}
          onOpenChange={setShowAddMemberDialog}
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
                <span className="font-medium text-mauve-12">
                  {memberToRemove?.name}
                </span>{" "}
                from this organization? They will lose access to all
                organization resources.
              </p>
            </div>
            {removeMemberMutation.isError && (
              <ErrorBanner text={removeMemberMutation.error.message} />
            )}
            <div className="flex justify-end gap-2">
              <Button
                type="button"
                intent="secondary"
                className="w-24"
                disabled={removeMemberMutation.isPending}
                onClick={() => setMemberToRemove(null)}
              >
                Cancel
              </Button>
              <Button
                className="w-24"
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

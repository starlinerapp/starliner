import React, { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { AvatarIcon } from "~/components/atoms/avatar/Avatar";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";

interface FormInput {
  email: string;
}

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

  const sendInviteMutation = useMutation(
    trpc.organization.sendInvite.mutationOptions(),
  );

  const removeMemberMutation = useMutation(
    trpc.organization.removeOrganizationMember.mutationOptions(),
  );

  const { register, handleSubmit, reset, watch } = useForm<FormInput>({
    defaultValues: { email: "" },
  });

  const emailInput = watch("email", "");

  function onInviteMember(data: FormInput) {
    sendInviteMutation.mutate(
      {
        organizationId: organization.id,
        toEmail: data.email,
        inviteUrlPrefix: `${window.location.origin}/organizations/invite/`,
      },
      {
        onSuccess: () => {
          reset();
          setShowAddMemberDialog(false);
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
        <div className="border-mauve-6 overflow-hidden rounded-md border text-sm shadow-xs">
          <table className="w-full border-collapse">
            <thead className="h-14">
              <tr className="border-mauve-6 bg-gray-2 border-b">
                <th className="text-mauve-12 w-1/2 px-4 py-3 text-left text-xs font-bold uppercase">
                  Member
                </th>
                <th className="text-mauve-12 w-1/2 px-4 py-3 text-left text-xs font-bold uppercase">
                  Role
                </th>
                {organization.isOwner && (
                  <th className="w-[20%] px-4 py-3">
                    <div className="flex justify-end">
                      <Button
                        className="w-28 text-xs"
                        intent="secondary"
                        onClick={() => setShowAddMemberDialog(true)}
                      >
                        Invite Member
                      </Button>
                    </div>
                  </th>
                )}
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                Array.from({ length: 2 }).map((_, i) => (
                  <tr
                    key={i}
                    className="border-mauve-6 border-b last:border-b-0"
                  >
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-8 w-8 rounded-full" />

                        <div className="flex flex-col gap-1">
                          <Skeleton className="h-4 w-24" />
                          <Skeleton className="h-4 w-36" />
                        </div>
                      </div>
                    </td>

                    <td className="px-4 py-3">
                      <Skeleton className="h-4 w-16" />
                    </td>
                    {organization.isOwner && <td />}
                  </tr>
                ))
              ) : members?.length === 0 ? (
                <tr>
                  <td
                    colSpan={organization.isOwner ? 3 : 2}
                    className="text-mauve-11 px-4 py-3 text-sm"
                  >
                    No members yet.
                  </td>
                </tr>
              ) : (
                members?.map((member) => (
                  <tr
                    key={member.id}
                    className="border-mauve-6 border-b last:border-b-0"
                  >
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <AvatarIcon
                          name={member.name}
                          profilePicture={member.avatarUrl}
                        />

                        <div className="flex flex-col">
                          <span className="text-mauve-12 font-medium">
                            {member.name}
                          </span>

                          <span className="text-mauve-11 text-sm">
                            {member.email}
                          </span>
                        </div>
                      </div>
                    </td>
                    <td className="text-mauve-11 px-4 py-3">
                      {member.is_owner ? "Owner" : "Member"}
                    </td>
                    {organization.isOwner && (
                      <td className="px-4 py-3">
                        <div className="flex justify-end">
                          {!member.is_owner &&
                            member.id !== Number(user?.user_id) && (
                              <Button
                                className="text-mauve-12 w-20 text-xs"
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
                      </td>
                    )}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {organization.isOwner && (
        <Dialog
          open={showAddMemberDialog}
          onOpenChange={setShowAddMemberDialog}
        >
          <DialogContent>
            <div className="flex flex-col gap-4">
              <div className="flex flex-col gap-2">
                <h1>Invite Member</h1>
                <p className="text-mauve-11 text-sm">
                  Invite members via email. They&apos;ll receive a link to join
                  your organization.
                </p>
              </div>
              <form
                className="flex flex-col gap-3"
                onSubmit={handleSubmit(onInviteMember)}
              >
                <input
                  type="email"
                  className="border-mauve-6 text-mauve-12 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                  placeholder="Email*"
                  {...register("email")}
                />
                <div className="flex justify-end gap-2">
                  <Button
                    type="button"
                    intent="secondary"
                    className="w-24"
                    onClick={() => {
                      setShowAddMemberDialog(false);
                      reset();
                      sendInviteMutation.reset();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    className="h-10 w-24"
                    type="submit"
                    disabled={!emailInput || sendInviteMutation.isPending}
                  >
                    Invite
                  </Button>
                </div>
              </form>
            </div>
          </DialogContent>
        </Dialog>
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

import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
} from "~/components/atoms/dialog/Dialog";

interface FormInput {
  email: string;
}

export default function Members() {
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const [showAddMemberDialog, setShowAddMemberDialog] = useState(false);

  const { data: members, isLoading } = useQuery(
    trpc.organization.getOrganizationMembers.queryOptions({
      id: organization.id,
    }),
  );

  const sendInviteMutation = useMutation(
    trpc.organization.sendInvite.mutationOptions(),
  );

  function onInviteMember(data: FormInput) {
    sendInviteMutation.mutate(
      {
        organizationId: organization.id,
        toEmail: data.email,
        inviteUrlPrefix: `${window.location.origin}/organizations/invite/`,
      },
      {
        onSuccess: () => reset(),
      },
    );
  }

  const { register, handleSubmit, reset, watch } = useForm<FormInput>({
    defaultValues: { email: "" },
  });
  const emailInput = watch("email", "");

  return (
    <div className="flex flex-col px-8 py-4">
      <div className="flex min-h-10 w-full items-center justify-between">
        <h1 className="text-xl font-bold">Members</h1>
        {organization.isOwner && (
          <Dialog
            open={showAddMemberDialog}
            onOpenChange={setShowAddMemberDialog}
          >
            <DialogTrigger>
              <Button className="w-32">Invite Member</Button>
            </DialogTrigger>
            <DialogContent>
              <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                  <h1>Invite Member</h1>
                  <p className="text-mauve-11 text-sm">
                    They will receive an email with a link to join your
                    organization.
                  </p>
                </div>
                <form
                  className="flex flex-col gap-3"
                  onSubmit={handleSubmit(onInviteMember)}
                >
                  <input
                    type="email"
                    className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
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
      </div>
      <div className="w-full pt-[42px] xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
            Organization Members
          </div>
          {isLoading ? (
            <>
              {Array.from({ length: 5 }).map((_, i) => (
                <div
                  key={i}
                  className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
                >
                  <Skeleton className="h-5 w-24" />
                  <Skeleton className="h-5 w-36" />
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
                <span>{member.name}</span>
                <span className="text-mauve-11">{member.email}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

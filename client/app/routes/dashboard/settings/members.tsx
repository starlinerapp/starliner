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
      },
      {
        onSuccess: () => reset(),
      },
    );
  }

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { isDirty },
  } = useForm<FormInput>({ defaultValues: { email: "" } });
  const emailInput = watch("email", "");

  return (
    <div className="flex flex-col gap-8 px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="pt-1 text-xl font-bold">Members</h1>
        <Dialog open={showAddMemberDialog} onOpenChange={setShowAddMemberDialog}>
          <DialogTrigger>
            <Button className="w-32">Invite Member</Button>
          </DialogTrigger>
          <DialogContent>
            <h2 className="text-mauve-12 mb-4 text-lg font-bold">
              Invite Member
            </h2>
            <form
              className="flex flex-col gap-3"
              onSubmit={handleSubmit(onInviteMember)}
            >
              <div className="flex items-center gap-2">
                <input
                  type="email"
                  className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
                  placeholder="Email*"
                  {...register("email")}
                />
                <Button
                  className="h-10 w-24 text-xs"
                  type="submit"
                  disabled={!emailInput}
                >
                  Invite
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>
      <div className="w-full xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-3 text-xs font-bold uppercase">
            <span>Organization Members</span>
          </div>
          {isLoading ? (
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
    </div>
  );
}

import React from "react";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";

interface FormInput {
  email: string;
}

interface AddMemberDialogProps {
  organizationId: number;
  teamId?: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export default function AddMemberDialog({
  organizationId,
  teamId,
  open,
  onOpenChange,
}: AddMemberDialogProps) {
  const trpc = useTRPC();

  const { register, handleSubmit, reset, watch } = useForm<FormInput>({
    defaultValues: { email: "" },
  });

  const emailInput = watch("email", "");

  const sendInviteMutation = useMutation(
    trpc.organization.sendInvite.mutationOptions(),
  );

  function onInviteMember(data: FormInput) {
    sendInviteMutation.mutate(
      {
        organizationId: organizationId,
        toEmail: data.email,
        inviteUrlPrefix: `${window.location.origin}/organizations/invite/`,
        ...(teamId != null ? { teamId } : {}),
      },
      {
        onSuccess: () => {
          reset();
          onOpenChange(false);
        },
      },
    );
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <h1>Invite Member</h1>
            <p className="text-mauve-11 text-sm">
              Invite members via email. They&apos;ll receive a link to join your
              organization.
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
                  onOpenChange(false);
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
  );
}

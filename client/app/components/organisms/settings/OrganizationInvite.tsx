import React, { useState } from "react";
import Button from "~/components/atoms/button/Button";
import { useMutation } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function OrganizationInvite() {
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const [inviteLink, setInviteLink] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const createInviteMutation = useMutation(
    trpc.organization.createInvite.mutationOptions(),
  );

  function handleCreateInvite() {
    createInviteMutation.mutate(
      { organizationId: organization.id },
      {
        onSuccess: (data) => {
          const link = `${window.location.origin}/organizations/invite/${data.id}`;
          setInviteLink(link);
          setCopied(false);
        },
      },
    );
  }

  function handleCopy() {
    if (!inviteLink) return;
    navigator.clipboard.writeText(inviteLink);
    setCopied(true);
  }

  if (!organization.isOwner) return null;

  return (
    <div className="w-full xl:w-3/5">
      <div className="border-mauve-6 rounded-md border-1 text-sm">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
          Invite Members
        </div>
        <div className="flex flex-col gap-3 px-4 py-3">
          <p className="text-mauve-11 text-xs">
            Generate an invite link to share with others. Links expire after 7
            days.
          </p>
          {inviteLink ? (
            <div className="flex items-center gap-2">
              <input
                className="border-mauve-6 text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
                value={inviteLink}
                readOnly
              />
              <Button className="h-9 w-24 text-xs" onClick={handleCopy}>
                {copied ? "Copied" : "Copy"}
              </Button>
              <Button
                className="h-9 w-24 text-xs"
                intent="secondary"
                onClick={handleCreateInvite}
                disabled={createInviteMutation.isPending}
              >
                New Link
              </Button>
            </div>
          ) : (
            <Button
              className="h-9 w-40 text-xs"
              onClick={handleCreateInvite}
              disabled={createInviteMutation.isPending}
            >
              {createInviteMutation.isPending
                ? "Generating..."
                : "Generate Invite Link"}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}

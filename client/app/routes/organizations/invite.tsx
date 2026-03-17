import React, { useState } from "react";
import { useNavigate, useParams } from "react-router";
import { ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useMutation } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

export default function AcceptInvite() {
  const trpc = useTRPC();
  const navigate = useNavigate();
  const { inviteId } = useParams<{ inviteId: string }>();
  const [error, setError] = useState<string | null>(null);

  const acceptInviteMutation = useMutation(
    trpc.organization.acceptInvite.mutationOptions(),
  );

  function handleAccept() {
    if (!inviteId) return;
    acceptInviteMutation.mutate(
      { inviteId },
      {
        onSuccess: () => {
          navigate("/");
        },
        onError: (err) => {
          setError(err.message);
        },
      },
    );
  }

  return (
    <div className="flex w-[500px] flex-col gap-4">
      <h1 className="text-xl font-medium">Join Organization</h1>
      <p className="text-mauve-11 text-sm">
        You&#39;ve been invited to join an organization. Click the button below
        to accept the invite and get started.
      </p>
      {error && <ErrorBanner text={error} />}
      <Button
        size="md"
        onClick={handleAccept}
        disabled={acceptInviteMutation.isPending}
      >
        {acceptInviteMutation.isPending ? "Joining..." : "Accept Invite"}{" "}
        <ChevronRight className="w-4 stroke-3" />
      </Button>
    </div>
  );
}

import React, { useState } from "react";
import type { Route } from "./+types/invite";
import { useNavigate, useParams } from "react-router";
import { ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { caller } from "~/utils/trpc/server";
import { getServerSession } from "~/utils/auth/server";

export async function loader(args: Route.LoaderArgs) {
  const c = await caller(args);

  const url = new URL(args.request.url);
  const inviteId = url.searchParams.get("invite_id");
  if (inviteId === null) throw new Error("Invalid invite");

  const session = await getServerSession(args.request);
  const invite = await c.organization.getInvite({
    inviteId,
  });

  if (invite.email !== session?.user.email) {
    throw new Error("Invalid invite");
  }
}

export default function AcceptInvite() {
  const trpc = useTRPC();
  const navigate = useNavigate();
  const { inviteId } = useParams<{ inviteId: string }>();
  const [error, setError] = useState<string | null>(null);

  const { data: invite, isLoading } = useQuery(
    trpc.organization.getInvite.queryOptions({
      inviteId: inviteId!,
    }),
  );

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
      {isLoading ? (
        <Skeleton className="h-5 w-64" />
      ) : invite ? (
        <p className="text-mauve-11 text-sm">
          You&#39;ve been invited to join{" "}
          <span className="text-mauve-12 font-medium">
            {invite.organization_name}
          </span>
          . Click the button below to accept the invite and get started.
        </p>
      ) : (
        <p className="text-mauve-11 text-sm">
          You&#39;ve been invited to join an organization. Click the button
          below to accept the invite and get started.
        </p>
      )}
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

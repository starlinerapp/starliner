import React, { useState } from "react";
import type { Route } from "./+types/invite";
import { useNavigate, useLoaderData } from "react-router";
import { ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useMutation } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { caller } from "~/utils/trpc/server";
import { getServerSession } from "~/utils/auth/server";
import { authClient } from "~/utils/auth/client";

export async function loader(args: Route.LoaderArgs) {
  const inviteId = args.params.inviteId;
  if (!inviteId) {
    throw new Response("Invalid invite", { status: 400 });
  }

  const session = await getServerSession(args.request);
  const c = await caller(args);
  const invite = await c.organization.getInvite({
    inviteId,
  });

  return {
    invite: {
      id: inviteId,
      organizationSlug: invite.organization_slug,
      organizationName: invite.organization_name,
    },
    emailMatches: invite.email === session?.user.email,
  };
}

export default function AcceptInvite() {
  const { invite, emailMatches } = useLoaderData<typeof loader>();

  const trpc = useTRPC();
  const navigate = useNavigate();

  const [error, setError] = useState<string | null>(
    emailMatches ? null : "You are not the intended recipient of this invite.",
  );

  const acceptInviteMutation = useMutation(
    trpc.organization.acceptInvite.mutationOptions(),
  );

  function handleAccept() {
    acceptInviteMutation.mutate(
      { inviteId: invite.id },
      {
        onSuccess: () => {
          navigate(`/${invite.organizationSlug}`);
        },
        onError: (err) => {
          setError(err.message);
        },
      },
    );
  }

  async function handleLogout() {
    await authClient.signOut();
    navigate(
      `/login?redirectTo=${encodeURIComponent(window.location.pathname)}`,
    );
  }

  return (
    <div className="flex w-125 flex-col gap-4">
      <h1 className="text-xl font-medium">Join {invite.organizationName}</h1>
      {error && (
        <ErrorBanner text={error}>
          {!emailMatches && (
            <button
              className="cursor-pointer text-sm font-light underline"
              type="button"
              onClick={handleLogout}
            >
              Logout
            </button>
          )}
        </ErrorBanner>
      )}
      <p className="text-mauve-11 text-sm">
        You&#39;ve been invited to join {invite.organizationName}. Click the
        button below to accept the invite and get started.
      </p>
      <Button
        size="md"
        onClick={handleAccept}
        disabled={acceptInviteMutation.isPending || !emailMatches}
      >
        {acceptInviteMutation.isPending ? "Joining..." : "Accept Invite"}{" "}
        <ChevronRight className="w-4 stroke-3" />
      </Button>
    </div>
  );
}

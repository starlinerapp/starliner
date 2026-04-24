import React from "react";
import { useLoaderData, useNavigate, useParams } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import TeamMembers from "~/components/organisms/settings/TeamMembers";
import { RepositoryAccess } from "~/components/organisms/settings/RepositoryAccess";
import { ClusterAccess } from "~/components/organisms/settings/ClusterAccess";

export function loader() {
  return {
    githubAppName: process.env.GITHUB_APP_NAME,
  };
}

export default function TeamDetail() {
  const { githubAppName } = useLoaderData<typeof loader>();
  const { teamId } = useParams();
  const navigate = useNavigate();

  const trpc = useTRPC();
  const organization = useOrganizationContext();

  const { data: teams, isLoading } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );
  const team = teams?.find((t) => t.id === Number(teamId));

  return (
    <div className="flex flex-col px-8 py-3">
      <div className="flex flex-col gap-1">
        <nav className="text-mauve-11 flex items-center gap-1 pb-4 text-sm">
          <button
            onClick={() => navigate("../", { relative: "path" })}
            className="hover:text-mauve-12 cursor-pointer hover:underline"
          >
            Teams
          </button>
          <span>&gt;</span>
          {isLoading ? (
            <Skeleton className="h-4 w-16" />
          ) : (
            <span className="text-mauve-12">#{team?.slug}</span>
          )}
        </nav>
        {isLoading ? (
          <Skeleton className="h-6 w-32" />
        ) : (
          <h1 className="text-xl font-bold">#{team?.slug}</h1>
        )}
      </div>
      <div className="flex flex-col gap-4">
        <div className="w-full pt-12">
          <TeamMembers teamId={Number(teamId)} />
        </div>

        {organization.isOwner && (
          <RepositoryAccess
            teamId={Number(teamId)}
            githubAppName={githubAppName}
          />
        )}

        {organization.isOwner && <ClusterAccess teamId={Number(teamId)} />}
      </div>
    </div>
  );
}

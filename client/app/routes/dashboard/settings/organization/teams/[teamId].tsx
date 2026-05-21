import React from "react";
import { useLoaderData, useNavigate, useParams } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import TeamMembers from "~/components/organisms/settings/TeamMembers";
import { RepositoryAccess } from "~/components/organisms/settings/RepositoryAccess";
import { ClusterAccess } from "~/components/organisms/settings/ClusterAccess";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";

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
    <>
      <Breadcrumbs
        crumbs={[
          { label: "Settings" },
          { label: "Organization" },
          {
            label: (
              <p
                className="cursor-pointer hover:underline"
                onClick={() => navigate("../", { relative: "path" })}
              >
                Teams
              </p>
            ),
          },
          {
            label: isLoading ? (
              <Skeleton className="h-4 w-16" />
            ) : (
              <p>#{team?.slug}</p>
            ),
          },
        ]}
      />
      <div className="flex flex-col gap-4 px-4 py-3">
        {isLoading ? (
          <Skeleton className="h-7 w-36" />
        ) : (
          <h1 className="text-xl font-bold">#{team?.slug}</h1>
        )}
        <div className="flex flex-col gap-4">
          <div className="w-full">
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
    </>
  );
}

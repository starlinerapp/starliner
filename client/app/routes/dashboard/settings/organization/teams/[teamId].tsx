import { useQuery } from "@tanstack/react-query";
import { useLoaderData, useNavigate, useParams } from "react-router";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import { ClusterAccess } from "~/components/organisms/settings/organization/team/ClusterAccess";
import { RepositoryAccess } from "~/components/organisms/settings/organization/team/RepositoryAccess";
import TeamDangerZone from "~/components/organisms/settings/organization/team/TeamDangerZone";
import TeamMembers from "~/components/organisms/settings/organization/team/TeamMembers";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

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
          <h1 className="font-bold text-xl">#{team?.slug}</h1>
        )}
        <div className="flex flex-col gap-4">
          <div className="w-full">
            <TeamMembers teamId={Number(teamId)} />
          </div>

          <RepositoryAccess
            teamId={Number(teamId)}
            githubAppName={githubAppName}
          />
          <ClusterAccess teamId={Number(teamId)} />
          {organization.isOwner && <TeamDangerZone teamId={Number(teamId)} />}
        </div>
      </div>
    </>
  );
}

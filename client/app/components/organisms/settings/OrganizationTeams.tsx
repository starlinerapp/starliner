import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { Link, useParams } from "react-router";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

export default function OrganizationTeams() {
  const trpc = useTRPC();
  const { slug } = useParams();
  const organization = useOrganizationContext();

  const { data: teamsData, isLoading: isTeamsLoading } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );

  const { data: members, isLoading: isMembersLoading } = useQuery(
    trpc.organization.getOrganizationMembers.queryOptions({
      id: organization.id,
    }),
  );

  return (
    <div className="w-full">
      <div className="border-mauve-6 w-full rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center justify-between border-b px-4 text-xs font-bold uppercase">
          <p>Your Teams</p>
        </div>
        {isTeamsLoading || isMembersLoading ? (
          <>
            {Array.from({ length: 3 }).map((_, i) => (
              <div
                key={i}
                className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
              >
                <Skeleton className="h-5 w-24" />
              </div>
            ))}
          </>
        ) : teamsData?.length === 0 ? (
          <div className="text-mauve-11 px-4 py-3 text-sm">No teams yet.</div>
        ) : (
          teamsData?.map((team) => (
            <Link
              key={team.id}
              to={`/${slug}/settings/teams/${team.id}`}
              className="border-mauve-6 text-mauve-12 hover:bg-gray-2 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
            >
              <div className="flex items-center gap-3">
                <div className="bg-violet-9 flex h-9 w-9 items-center justify-center rounded-md text-base text-white">
                  {team.slug.substring(0, 1)?.toUpperCase()}
                </div>
                <div className="flex flex-col">
                  <span className="text-violet-11">#{team.slug}</span>
                  <p className="text-mauve-11 text-sm">
                    {members?.length ?? 0}{" "}
                    {(members?.length ?? 0) === 1 ? "Member" : "Members"}
                  </p>
                </div>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}

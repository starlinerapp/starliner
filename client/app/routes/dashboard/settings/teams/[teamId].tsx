import React from "react";
import { useParams, useNavigate } from "react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Button from "~/components/atoms/button/Button";

export default function TeamDetail() {
  const { teamId } = useParams();
  const trpc = useTRPC();
  const organization = useOrganizationContext();

  const queryClient = useQueryClient();

  const { data: members, isLoading } = useQuery(
    trpc.team.getTeamMembers.queryOptions({
      organizationId: organization.id,
      teamId: Number(teamId),
    }),
  );

  const leaveTeamMutation = useMutation(
    trpc.team.removeTeamMember.mutationOptions(),
  );

  return (
    <div className="flex flex-col gap-8 px-8 py-4">
      <h1 className="pt-1 text-xl font-bold">Team Members</h1>
      <div className="w-full xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
            Members
          </div>
          {isLoading ? (
            <div className="text-mauve-11 px-4 py-3 text-sm">Loading...</div>
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
                <span>{member.better_auth_id}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

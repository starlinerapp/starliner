import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

export default function Members() {
  const trpc = useTRPC();
  const organization = useOrganizationContext();

  const { data: members, isLoading } = useQuery(
    trpc.organization.getOrganizationMembers.queryOptions({
      id: organization.id,
    }),
  );

  return (
    <div className="flex flex-col px-8 py-4">
      <h1 className="pt-1 text-xl font-bold">Members</h1>
      <div className="w-full pt-10 xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
            Organization Members
          </div>
          {isLoading ? (
            <>
              {Array.from({ length: 5 }).map((_, i) => (
                <div
                  key={i}
                  className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
                >
                  <Skeleton className="h-5 w-24" />
                  <Skeleton className="h-5 w-36" />
                </div>
              ))}
            </>
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
                <span>{member.name}</span>
                <span className="text-mauve-11 text-xs">{member.email}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

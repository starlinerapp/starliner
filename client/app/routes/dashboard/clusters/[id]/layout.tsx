import { Outlet, useParams } from "react-router";
import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function ClusterLayout() {
  const { slug, id } = useParams<{
    slug: string;
    id: string;
  }>();

  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const { data: clusters, isLoading } = useQuery(
    trpc.organization.getOrganizationClusters.queryOptions({
      id: Number(organization.id),
    }),
  );

  const currentCluster = clusters?.find((cluster) => cluster.id === Number(id));

  const navigationBarItems = [
    {
      title: "General",
      href: `/${slug}/clusters/${id}/general`,
    },
    {
      title: "Settings",
      href: `/${slug}/clusters/${id}/settings`,
    },
  ];

  return (
    <div className="bg-violet-1 flex h-full flex-col">
      {isLoading ? (
        <div className="px-4 pt-4">
          <Skeleton className="h-7 w-32" />
        </div>
      ) : (
        <div className="flex items-center gap-3 px-4 pt-4">
          <h1 className="text-mauve-12 text-xl font-bold">
            {currentCluster?.name}
          </h1>
        </div>
      )}
      <LinkNavigationBar items={navigationBarItems} />
      <Outlet />
    </div>
  );
}

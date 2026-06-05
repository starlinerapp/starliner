import { useQuery } from "@tanstack/react-query";
import { Outlet, redirect, useParams } from "react-router";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";
import { caller } from "~/utils/trpc/server";
import type { Route } from "./+types/layout";

export async function loader(loaderArgs: Route.LoaderArgs) {
  const { params } = loaderArgs;

  const trpcCaller = await caller(loaderArgs);

  try {
    await trpcCaller.cluster.getCluster({
      id: Number(params.id),
    });
  } catch {
    return redirect(`/${params.slug}/clusters/all`);
  }
}

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
    <div className="flex h-full flex-col bg-violet-1">
      {isLoading ? (
        <div className="bg-violet-1 px-4 pt-4">
          <Skeleton className="h-7 w-32" />
        </div>
      ) : (
        <div className="flex items-center gap-3 bg-violet-1 px-4 pt-4">
          <h1 className="font-bold text-mauve-12 text-xl">
            {currentCluster?.name}
          </h1>
        </div>
      )}
      <div className="bg-mauve-1">
        <LinkNavigationBar items={navigationBarItems} />
      </div>
      <div className="h-full w-full bg-white">
        <Outlet />
      </div>
    </div>
  );
}

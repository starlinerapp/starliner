import React from "react";
import type { Route } from "./+types/layout";
import {
  Outlet,
  redirect,
  useLoaderData,
  useLocation,
  useParams,
} from "react-router";
import { caller } from "~/utils/trpc/server";
import { auth } from "~/utils/auth/server";
import Sidebar from "~/components/organisms/sidebar/Sidebar";
import ExtendedSidebar from "~/components/organisms/extended-sidebar/ExtendedSidebar";
import { Cog, InboxStack, Servers } from "~/components/atoms/icons";
import { OrganizationProvider } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";

export async function loader(loaderArgs: Route.LoaderArgs) {
  const { request, params } = loaderArgs;
  const session = await auth.api.getSession({
    headers: request.headers,
  });

  if (!session) {
    return redirect("/login");
  }

  const trpc = await caller(loaderArgs);
  const organizations = await trpc.organization.getUserOrganizations();

  if (!organizations || organizations.length === 0) {
    return redirect("/organizations/new");
  }

  if (!params.slug) {
    return redirect(`/${organizations[0].slug}/projects`);
  }

  const userOrganization = organizations.find((o) => o.slug === params.slug);
  if (!userOrganization) {
    throw new Response(undefined, { status: 404 });
  }

  return {
    organization: userOrganization,
  };
}

export default function Layout() {
  const { organization } = useLoaderData<typeof loader>();

  const location = useLocation();
  const { slug } = useParams<{ slug: string }>();

  const trpc = useTRPC();
  const { data: projects, isLoading: isProjectsLoading } = useQuery(
    trpc.organization.getOrganizationProjects.queryOptions({
      id: organization.id,
    }),
  );

  const { data: clusters, isLoading: isClustersLoading } = useQuery(
    trpc.organization.getOrganizationClusters.queryOptions({
      id: organization.id,
    }),
  );

  const sidebarItems = [
    {
      id: "projects",
      title: "Projects",
      icon: <InboxStack />,
      href: `/${slug}/projects`,
      extended: [
        [
          {
            id: "all-projects",
            title: "All Projects",
            href: `/${slug}/projects/all`,
          },
        ],
        [
          ...(projects?.map((project) => ({
            id: `project-${project.id}`,
            title: project.name ?? "",
            href: `/${slug}/projects/${project.id}`,
          })) ?? []),
        ],
      ],
    },
    {
      id: "clusters",
      title: "Clusters",
      icon: <Servers />,
      href: `/${slug}/clusters`,
      extended: [
        [
          {
            id: "all-clusters",
            title: "All Clusters",
            href: `/${slug}/clusters/all`,
          },
        ],
        [
          ...(clusters?.map((cluster) => ({
            id: `cluster-${cluster.id}`,
            title: cluster.name ?? "",
            href: `/${slug}/clusters/${cluster.id}`,
          })) ?? []),
        ],
      ],
    },
    {
      id: "settings",
      title: "Settings",
      icon: <Cog />,
      href: `/${slug}/settings`,
    },
  ];

  const activeItem = sidebarItems.find((item) =>
    location.pathname.startsWith(item.href),
  );

  return (
    <OrganizationProvider
      name={organization.name}
      id={organization.id}
      slug={organization.slug}
    >
      <Sidebar sidebarItems={sidebarItems}>
        <ExtendedSidebar
          title={activeItem?.title ?? ""}
          sections={activeItem?.extended ?? []}
          isLoading={isProjectsLoading || isClustersLoading}
        >
          <Outlet />
        </ExtendedSidebar>
      </Sidebar>
    </OrganizationProvider>
  );
}

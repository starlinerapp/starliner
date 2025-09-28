import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";
import { Outlet, useParams } from "react-router";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function ProjectLayout() {
  const { slug, id } = useParams<{ slug: string; id: string }>();
  const organization = useOrganizationContext();
  const trpc = useTRPC();

  const { data: projects, isLoading } = useQuery(
    trpc.organization.getOrganizationProjects.queryOptions({
      id: organization.id,
    }),
  );

  const currentProject = projects?.find((p) => p.id === Number(id));

  const navigationBarItems = [
    { title: "Architecture", href: `/${slug}/projects/${id}/architecture` },
    { title: "Observability", href: `/${slug}/projects/${id}/observability` },
    { title: "Logs", href: `/${slug}/projects/${id}/logs` },
    { title: "Settings", href: `/${slug}/projects/${id}/settings` },
  ];

  return (
    <div className="flex h-full flex-col">
      <NavigationBar
        isLoading={isLoading}
        title={currentProject?.name ?? ""}
        items={navigationBarItems}
      />
      <Outlet />
    </div>
  );
}

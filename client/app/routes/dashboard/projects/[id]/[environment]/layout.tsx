import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";
import { Outlet, useParams } from "react-router";
import React, { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

export default function ProjectLayout() {
  const { slug, id, environment } = useParams<{
    slug: string;
    id: string;
    environment: string;
  }>();
  const organization = useOrganizationContext();
  const trpc = useTRPC();

  const { data: projects, isLoading } = useQuery(
    trpc.organization.getOrganizationProjects.queryOptions({
      id: organization.id,
    }),
  );

  const currentProject = projects?.find((p) => p.id === Number(id));

  const navigationBarItems = [
    {
      title: "Architecture",
      href: `/${slug}/projects/${id}/${environment}/architecture`,
    },
    {
      title: "Builds",
      href: `/${slug}/projects/${id}/${environment}/builds`,
    },
    {
      title: "Settings",
      href: `/${slug}/projects/${id}/${environment}/settings`,
    },
  ];

  const environments = currentProject?.environments.map((env) => env) ?? [];

  const [selectedEnvironment, setSelectedEnvironment] = React.useState<
    string | undefined
  >(environment);

  useEffect(() => {
    if (environments.length > 0 && !selectedEnvironment)
      setSelectedEnvironment(environments[0].slug);
  }, [environments]);

  return (
    <div className="flex h-full flex-col">
      <div className="bg-violet-1">
        {isLoading ? (
          <div className="px-4 pt-4">
            <Skeleton className="h-7 w-32" />
          </div>
        ) : (
          <div className="flex items-center gap-3 px-4 pt-4">
            <h1 className="text-mauve-12 text-xl font-bold">
              {currentProject?.name}
            </h1>
            <div className="border-violet-10 flex items-center gap-1.5 rounded-md border-[1px] px-2 text-sm">
              <h1>
                {environments.find((e) => e.slug === selectedEnvironment)?.name}
              </h1>
            </div>
          </div>
        )}
        <LinkNavigationBar items={navigationBarItems} />
      </div>
      <div className="h-[calc(100vh-90px)] overflow-y-scroll">
        <Outlet />
      </div>
    </div>
  );
}

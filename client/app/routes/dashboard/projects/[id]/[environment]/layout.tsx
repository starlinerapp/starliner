import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";
import { Outlet, useLocation, useParams } from "react-router";
import React, { useEffect, useMemo } from "react";
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
  const location = useLocation();

  const organization = useOrganizationContext();
  const trpc = useTRPC();

  const { data: projects, isLoading: isProjectsLoading } = useQuery(
    trpc.organization.getOrganizationProjects.queryOptions({
      id: organization.id,
    }),
  );

  const currentProject = projects?.find((p) => p.id === Number(id));

  const currentEnvironment = useMemo(
    () => currentProject?.environments.find((e) => e.slug === environment),
    [currentProject, environment],
  );

  const { data: environmentBuilds, isLoading: isEnvironmentBuildsLoading } =
    useQuery({
      ...trpc.environment.getEnvironmentBuilds.queryOptions({
        id: Number(currentEnvironment?.id),
      }),
      enabled: !!currentEnvironment,
      refetchOnWindowFocus: "always",
      refetchOnMount: "always",
    });

  const [initialBuildCount, setInitialBuildCount] = React.useState<
    number | null
  >(null);

  useEffect(() => {
    if (!currentEnvironment) {
      setInitialBuildCount(null);
      return;
    }

    if (initialBuildCount === null && environmentBuilds !== undefined) {
      setInitialBuildCount(environmentBuilds.length);
    }
  }, [currentEnvironment?.id, environmentBuilds, initialBuildCount]);

  useEffect(() => {
    if (!environmentBuilds) return;

    if (location.pathname.endsWith("/builds")) {
      setTimeout(() => {
        setInitialBuildCount(environmentBuilds.length);
      }, 1000);
    }
  }, [location.pathname, environmentBuilds]);

  const currentBuildCount = environmentBuilds?.length ?? 0;

  const newBuildsSinceFirstLoad =
    initialBuildCount === null
      ? 0
      : Math.max(0, currentBuildCount - initialBuildCount);

  const navigationBarItems = [
    {
      title: "Architecture",
      href: `/${slug}/projects/${id}/${environment}/architecture`,
    },
    {
      title: (
        <span className="flex items-center gap-2">
          <p>Builds</p>
          {newBuildsSinceFirstLoad > 0 && !isEnvironmentBuildsLoading && (
            <span className="bg-violet-9 rounded-full px-2 py-0.5 text-xs text-white">
              + {newBuildsSinceFirstLoad}
            </span>
          )}
        </span>
      ),
      href: `/${slug}/projects/${id}/${environment}/builds`,
    },
    {
      title: "Settings",
      href: `/${slug}/projects/${id}/${environment}/settings`,
    },
  ];

  return (
    <div className="flex h-full flex-col">
      <div className="bg-violet-1">
        {isProjectsLoading || isEnvironmentBuildsLoading ? (
          <div className="px-4 pt-4">
            <Skeleton className="h-7 w-32" />
          </div>
        ) : (
          <div className="flex items-center gap-3 px-4 pt-4">
            <h1 className="text-mauve-12 text-xl font-bold">
              {currentProject?.name}
            </h1>
            <div className="border-violet-10 flex items-center gap-1.5 rounded-md border-[1px] px-2 text-sm">
              <h1>{currentEnvironment?.name}</h1>
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

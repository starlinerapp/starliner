import React, { useMemo } from "react";
import BuildCard from "~/components/organisms/build-card/BuildCard";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useParams } from "react-router";

export default function Builds() {
  const trpc = useTRPC();
  const { id, environment } = useParams<{
    id: string;
    environment: string;
  }>();

  const { data: project } = useQuery(
    trpc.project.getProject.queryOptions({ id: Number(id) }),
  );

  const currentEnvironment = useMemo(
    () => project?.environments.find((e) => e.slug === environment),
    [project, environment],
  );

  const { data: environmentBuilds, isLoading } = useQuery({
    ...trpc.environment.getEnvironmentBuilds.queryOptions({
      id: Number(currentEnvironment?.id),
    }),
    enabled: !!currentEnvironment,
    refetchOnWindowFocus: "always",
    refetchOnMount: "always",
    refetchInterval: (query) => {
      const builds = query.state.data;
      if (!builds) return false;
      const shouldPoll = builds?.some(
        (build) => build.status === "building" || build.status === "queued",
      );
      return shouldPoll ? 1000 : false;
    },
  });

  if (!isLoading && environmentBuilds?.length === 0)
    return (
      <div className="flex flex-col gap-1 p-4">
        <p className="text-mauve-11">
          There are no builds for this environment yet.
        </p>
      </div>
    );

  return (
    <div className="flex flex-col gap-4 p-4">
      {environmentBuilds?.map((build, i) => (
        <BuildCard
          isCollapsed={i > 0}
          key={i}
          buildId={build.buildId}
          commitHash={build.commitHash}
          source={build.source}
          serviceName={build.deploymentName}
          createdAt={build.createdAt}
          status={build.status}
        />
      ))}
    </div>
  );
}

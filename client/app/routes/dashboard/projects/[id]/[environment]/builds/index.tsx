import React, { useMemo } from "react";
import LogsCard from "~/components/organisms/logs-card/LogsCard";
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

  const { data: environmentBuilds } = useQuery({
    ...trpc.environment.getEnvironmentBuilds.queryOptions({
      id: Number(currentEnvironment?.id),
    }),
    enabled: !!currentEnvironment,
    refetchOnWindowFocus: true,
    refetchInterval: (query) => {
      const builds = query.state.data;
      if (!builds) return false;
      const shouldPoll = builds?.some(
        (build) => build.status === "building" || build.status === "queued",
      );
      return shouldPoll ? 1000 : false;
    },
  });

  return (
    <div className="flex flex-col gap-4 p-4">
      {environmentBuilds?.map((build, i) => (
        <LogsCard
          key={i}
          buildId={build.buildId}
          serviceName={build.deploymentName}
          createdAt={build.createdAt}
          status={build.status}
        />
      ))}
    </div>
  );
}

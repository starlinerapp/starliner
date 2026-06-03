import React, { useMemo } from "react";
import DeploymentCard from "~/components/organisms/deployment-card/DeploymentCard";
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
        (build) =>
          build.status === "building" ||
          build.status === "queued" ||
          (build.status === "success" &&
            build.deploymentRolloutStatus === "pending"),
      );
      return shouldPoll ? 1000 : false;
    },
  });

  const { data: environmentDeployments } = useQuery({
    ...trpc.environment.getEnvironmentDeployments.queryOptions({
      id: Number(currentEnvironment?.id),
    }),
    enabled: !!currentEnvironment,
  });

  const skippedBuildDeploymentIds = useMemo(() => {
    if (!environmentDeployments) {
      return new Set<number>();
    }

    return new Set([
      ...environmentDeployments.images.map((d) => d.id),
      ...environmentDeployments.databases.map((d) => d.id),
      ...environmentDeployments.ingresses.map((d) => d.id),
    ]);
  }, [environmentDeployments]);

  if (!isLoading && environmentBuilds?.length === 0)
    return (
      <div className="flex flex-col gap-1 p-4">
        <p className="text-mauve-11">
          There are no deployments for this environment yet.
        </p>
      </div>
    );

  return (
    <div className="flex flex-col gap-4 p-4">
      {environmentBuilds?.map((build, i) => (
        <DeploymentCard
          isCollapsed={i > 0}
          key={build.buildId}
          buildId={build.buildId}
          deploymentId={build.deploymentId}
          commitHash={build.commitHash}
          source={build.source}
          serviceName={build.deploymentName}
          createdAt={build.createdAt}
          status={build.status}
          deploymentRolloutStatus={build.deploymentRolloutStatus}
          isDeployOnly={skippedBuildDeploymentIds.has(build.deploymentId)}
        />
      ))}
    </div>
  );
}

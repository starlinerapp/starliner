import React, { useMemo } from "react";
import BuildCard from "~/components/organisms/build-card/BuildCard";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useNavigate, useParams } from "react-router";
import { Box } from "lucide-react";
import Button from "~/components/atoms/button/Button";

export default function Builds() {
  const navigate = useNavigate();
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
      <div className="flex h-full flex-col items-center justify-center gap-4 p-4">
        <div className="bg-mauve-3 relative flex h-12 w-12 items-center justify-center rounded-sm">
          {/* top left */}
          <div className="border-mauve-8 absolute top-0 left-0 h-3 w-3 rounded-tl-sm border-t-2 border-l-2" />

          {/* top right */}
          <div className="border-mauve-8 absolute top-0 right-0 h-3 w-3 rounded-tr-sm border-t-2 border-r-2" />

          {/* bottom left */}
          <div className="border-mauve-8 absolute bottom-0 left-0 h-3 w-3 rounded-bl-sm border-b-2 border-l-2" />

          {/* bottom right */}
          <div className="border-mauve-8 absolute right-0 bottom-0 h-3 w-3 rounded-br-sm border-r-2 border-b-2" />

          <Box className="text-mauve-9 h-8 w-8" />
        </div>
        <div className="flex flex-col items-center">
          <p className="text-mauve-12 text-sm">
            There are no builds in this environment yet.
          </p>
          <p className="text-mauve-11 text-sm">
            Deploy your first service and view the build logs here.
          </p>
        </div>
        <Button
          onClick={() => navigate("../architecture", { relative: "path" })}
          className="w-32"
        >
          Deploy Service
        </Button>
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

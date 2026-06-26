import { useQuery } from "@tanstack/react-query";
import { useSubscription } from "@trpc/tanstack-react-query";
import { Play } from "lucide-react";
import { useEffect, useState } from "react";
import LogsConsole from "~/components/molecules/logs-console/LogsConsole";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

interface DeploymentTabProps {
  isActive: boolean;
  onSelect: () => void;
}

export function DeploymentTab({ isActive, onSelect }: DeploymentTabProps) {
  return (
    <div className="relative">
      <div className="absolute top-1/2 -left-1 h-2 w-2 -translate-y-1/2 rounded-full bg-mauve-8" />
      <button
        type="button"
        onClick={onSelect}
        className={cn(
          "relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border bg-white px-4 py-0.5 hover:bg-mauve-2",
          isActive
            ? "border-violet-9 bg-violet-3 text-violet-9"
            : "border-mauve-6 text-mauve-9",
        )}
      >
        <div
          className={cn(
            "flex rounded-full border-[1.5px] p-0.5",
            isActive ? "border-violet-9" : "border-mauve-9",
          )}
        >
          <Play
            className={cn(
              "h-2 w-2",
              isActive
                ? "fill-violet-9 stroke-violet-9"
                : "fill-mauve-9 stroke-mauve-9",
            )}
          />
        </div>
        Deploy
      </button>
    </div>
  );
}

interface DeploymentLogsProps {
  deploymentId: number;
  buildStatus: string;
  deploymentRolloutStatus: string;
  isDeployOnly?: boolean;
  enabled?: boolean;
}

export function DeploymentLogs({
  deploymentId,
  buildStatus,
  deploymentRolloutStatus,
  isDeployOnly = false,
  enabled = true,
}: DeploymentLogsProps) {
  const trpc = useTRPC();
  const [lines, setLines] = useState<string[]>([]);

  useEffect(() => {
    setLines([]);
  }, [deploymentId]);

  const buildComplete = buildStatus === "success";
  const buildFailed = buildStatus === "failure";
  const isDeploying = isDeployOnly
    ? deploymentRolloutStatus === "pending"
    : buildComplete && deploymentRolloutStatus === "pending";
  const isDeployComplete =
    deploymentRolloutStatus === "success" ||
    deploymentRolloutStatus === "failure";

  useSubscription(
    trpc.deployment.streamDeploymentStatusLogs.subscriptionOptions(
      { deploymentId },
      {
        enabled: enabled && isDeploying,
        onData: (chunk) => {
          const line = chunk.replace(/\r$/, "");
          if (!line) {
            return;
          }

          setLines((prev) => [...prev, line]);
        },
      },
    ),
  );

  const { data: completedLogLines } = useQuery({
    ...trpc.deployment.getDeploymentStatusLogs.queryOptions({ deploymentId }),
    enabled: enabled && buildComplete && !isDeploying && isDeployComplete,
  });

  const displayedLines = isDeploying ? lines : (completedLogLines ?? []);

  if (buildFailed) {
    return (
      <pre className="whitespace-pre-wrap text-mauve-11 text-sm">
        Build failed — deployment was not triggered.
      </pre>
    );
  }

  return <LogsConsole logs={displayedLines} resetKey={deploymentId} />;
}

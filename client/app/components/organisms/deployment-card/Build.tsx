import { useQuery } from "@tanstack/react-query";
import { useSubscription } from "@trpc/tanstack-react-query";
import { Hammer } from "lucide-react";
import { useEffect, useState } from "react";
import LogsConsole from "~/components/molecules/logs-console/LogsConsole";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

interface BuildTabProps {
  isActive: boolean;
  onSelect: () => void;
}

export function BuildTab({ isActive, onSelect }: BuildTabProps) {
  return (
    <div className="relative">
      <div className="absolute top-1/2 -right-1 h-2 w-2 -translate-y-1/2 rounded-full bg-mauve-8" />
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
          <Hammer
            className={cn(
              "h-2 w-2",
              isActive
                ? "fill-violet-9 stroke-violet-9"
                : "fill-mauve-9 stroke-mauve-9",
            )}
          />
        </div>
        Build
      </button>
    </div>
  );
}

interface BuildLogsProps {
  buildId: number;
  buildStatus: string;
  enabled?: boolean;
}

export function BuildLogs({
  buildId,
  buildStatus,
  enabled = true,
}: BuildLogsProps) {
  const trpc = useTRPC();
  const [logs, setLogs] = useState<string[]>([]);

  const isLive = buildStatus === "queued" || buildStatus === "building";
  const isComplete = buildStatus === "success" || buildStatus === "failure";

  useEffect(() => {
    setLogs([]);
  }, [buildId]);

  useSubscription(
    trpc.build.streamBuildLogs.subscriptionOptions(
      { buildId },
      {
        enabled: enabled && isLive,
        onData: (chunk) => {
          setLogs((prev) => [...prev, chunk]);
        },
      },
    ),
  );

  const { data: completedLogs } = useQuery({
    ...trpc.build.getBuildLogs.queryOptions({ id: buildId }),
    enabled: enabled && isComplete,
  });

  const displayedLogs = isLive
    ? logs
    : completedLogs?.logs
      ? completedLogs.logs.split("\n")
      : [];

  return <LogsConsole logs={displayedLogs} resetKey={buildId} />;
}

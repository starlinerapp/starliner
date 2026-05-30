import React, { useEffect, useRef, useState } from "react";
import { Play } from "lucide-react";
import { useSubscription } from "@trpc/tanstack-react-query";
import { useTRPC } from "~/utils/trpc/react";

const deploymentStatusSnapshotDelimiter =
  "────────────────────────────────────────";

function isSnapshotDelimiter(line: string) {
  return line === deploymentStatusSnapshotDelimiter || /^─{10,}$/.test(line);
}

interface DeploymentTabProps {
  isActive: boolean;
  onSelect: () => void;
}

export function DeploymentTab({ isActive, onSelect }: DeploymentTabProps) {
  return (
    <div className="relative">
      <div className="bg-mauve-8 absolute top-1/2 -left-1 h-2 w-2 -translate-y-1/2 rounded-full" />
      <button
        type="button"
        onClick={onSelect}
        className={
          isActive
            ? "border-violet-9 bg-violet-3 hover:bg-mauve-2 text-violet-9 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border px-4 py-0.5"
            : "border-mauve-6 hover:bg-mauve-2 text-mauve-9 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border bg-white px-4 py-0.5"
        }
      >
        <div
          className={
            isActive
              ? "border-violet-9 flex rounded-full border-[1.5px] p-0.5"
              : "border-mauve-9 flex rounded-full border-[1.5px] p-0.5"
          }
        >
          <Play
            className={
              isActive
                ? "fill-violet-9 stroke-violet-9 h-2 w-2"
                : "fill-mauve-9 stroke-mauve-9 h-2 w-2"
            }
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
  deploymentDeleted?: boolean;
}

export function DeploymentLogs({
  deploymentId,
  buildStatus,
  deploymentDeleted = false,
}: DeploymentLogsProps) {
  const trpc = useTRPC();
  const [logs, setLogs] = useState<string[]>([]);
  const logsScrollRef = useRef<HTMLPreElement>(null);
  const hasLoadedInitial = useRef(false);

  useEffect(() => {
    setLogs([]);
    hasLoadedInitial.current = false;
  }, [deploymentId]);

  const buildComplete = buildStatus === "success";
  const buildFailed = buildStatus === "failure";

  useSubscription(
    trpc.deployment.streamDeploymentStatusLogs.subscriptionOptions(
      { deploymentId },
      {
        enabled: buildComplete,
        onData: (chunk) => {
          setLogs((prev) => [...prev, chunk]);
        },
      },
    ),
  );

  useEffect(() => {
    const el = logsScrollRef.current;
    if (!el) {
      return;
    }
    const scrollToBottom = (behavior: ScrollBehavior) => {
      const top = el.scrollHeight - el.clientHeight;
      if (top <= 0) {
        return;
      }
      el.scrollTo({ top, left: 0, behavior });
    };
    if (!hasLoadedInitial.current) {
      if (logs.length > 0) {
        hasLoadedInitial.current = true;
        requestAnimationFrame(() => {
          requestAnimationFrame(() => {
            scrollToBottom("auto");
          });
        });
      }
      return;
    }
    scrollToBottom("smooth");
  }, [logs]);

  if (buildFailed) {
    return (
      <pre className="text-mauve-11 max-h-125 overflow-y-auto whitespace-pre-wrap">
        Build failed — deployment was not triggered.
      </pre>
    );
  }

  if (!buildComplete) {
    return (
      <pre className="text-mauve-11 max-h-125 overflow-y-auto whitespace-pre-wrap">
        Deployment will begin after the build completes.
      </pre>
    );
  }

  return (
    <pre
      ref={logsScrollRef}
      className="text-mauve-11 max-h-125 w-full overflow-y-auto font-mono text-sm break-all whitespace-pre-wrap"
    >
      {logs.map((line, i) =>
        isSnapshotDelimiter(line) ? (
          <hr key={i} className="border-mauve-6 my-3" />
        ) : (
          <span key={i} className="block">
            {line}
          </span>
        ),
      )}
    </pre>
  );
}

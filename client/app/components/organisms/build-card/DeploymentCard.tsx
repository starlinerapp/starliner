import React, { useEffect, useRef, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronRight } from "~/components/atoms/icons";
import { formatDistanceToNow } from "date-fns";
import { Spinner } from "~/components/atoms/spinner/Spinner";
import { Check, GitMerge, Hammer, Play, X } from "lucide-react";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useTRPC } from "~/utils/trpc/react";
import { useSubscription } from "@trpc/tanstack-react-query";

interface LogsCardProps {
  isCollapsed?: boolean;
  buildId: number;
  commitHash: string;
  source: string;
  serviceName: string;
  createdAt: string;
  status: string;
  args?: { name: string; value: string }[];
}

export default function DeploymentCard({
  isCollapsed: collapsed = true,
  buildId,
  commitHash,
  source,
  serviceName,
  status,
  createdAt,
}: LogsCardProps) {
  const trpc = useTRPC();

  const [isCollapsed, setIsCollapsed] = useState(collapsed);
  const [logs, setLogs] = useState<string[]>([]);
  const [hasReceivedLogs, setHasReceivedLogs] = useState(false);

  const logsScrollRef = useRef<HTMLPreElement>(null);
  const hasLoadedInitial = useRef(false);

  useEffect(() => {
    if (isCollapsed) {
      return;
    }
    setLogs([]);
    setHasReceivedLogs(false);
    hasLoadedInitial.current = false;
  }, [isCollapsed, buildId]);

  useSubscription(
    trpc.build.streamBuildLogs.subscriptionOptions(
      { buildId },
      {
        enabled: !isCollapsed,
        onData: (chunk) => {
          setLogs((prev) => [...prev, chunk]);
          setHasReceivedLogs(true);
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

  const isBuilding = status === "queued" || status === "building";

  return (
    <div className="shadow-xs">
      <div className="border-mauve-6 rounded-t-md border px-4 py-3 text-sm">
        <div className="flex gap-3">
          <div className="flex h-5 w-5 shrink-0 items-center justify-center">
            {isBuilding && <Spinner className="stroke-violet-10 size-5" />}
            {status === "success" && (
              <div className="bg-grass-9 flex h-4.5 w-4.5 items-center justify-center rounded-full">
                <Check className="w-3.5 stroke-white stroke-2" />
              </div>
            )}
            {status === "failure" && (
              <div className="bg-red-9 flex h-4.5 w-4.5 items-center justify-center rounded-full">
                <X className="w-3.5 stroke-white stroke-2" />
              </div>
            )}
          </div>

          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <span className="flex items-center gap-2">
                <p>{serviceName}</p>
              </span>
              <p>·</p>
              <p className="text-mauve-10">
                {formatDistanceToNow(new Date(createdAt))} ago
              </p>
              {commitHash && (
                <p className="text-mauve-10 bg-gray-3 border-mauve-6 flex items-center gap-1 rounded-md border px-1.5">
                  <GitMerge size={16} />
                  {commitHash.slice(0, 7)}
                </p>
              )}
            </div>

            <div className="text-mauve-10 mt-0.5 text-sm">
              <p>{source === "manual" ? "Manually triggered" : "On Push"}</p>
            </div>
          </div>
        </div>
      </div>
      <div className="border-mauve-6 rounded-b-md border-x border-b text-sm">
        <div
          className="flex cursor-pointer items-center gap-3 px-4 py-2 text-sm"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          <motion.div
            animate={{ rotate: isCollapsed ? 0 : 90 }}
            transition={{ duration: 0.2, ease: "easeOut" }}
          >
            <ChevronRight className="w-4 stroke-2" />
          </motion.div>
          <div
            className="relative flex items-center"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="relative">
              <div className="bg-mauve-8 absolute top-1/2 -right-1 h-2 w-2 -translate-y-1/2 rounded-full" />

              <span className="border-violet-9 bg-violet-3 hover:bg-mauve-2 text-violet-9 relative z-10 flex items-center gap-1.5 rounded-md border px-4 py-0.5">
                <div className="border-violet-9 flex rounded-full border-[1.5px] p-0.5">
                  <Hammer className="fill-violet-9 stroke-violet-9 h-2 w-2" />
                </div>
                Build
              </span>
            </div>

            <div className="bg-mauve-8 h-px w-4" />

            <div className="relative">
              <div className="bg-mauve-8 absolute top-1/2 -left-1 h-2 w-2 -translate-y-1/2 rounded-full" />

              <span className="border-mauve-6 hover:bg-mauve-2 text-mauve-9 relative z-10 flex items-center gap-1.5 rounded-md border bg-white px-4 py-0.5">
                <div className="border-mauve-9 flex rounded-full border-[1.5px] p-0.5">
                  <Play className="fill-mauve-9 stroke-mauve-9 h-2 w-2" />
                </div>
                Deploy
              </span>
            </div>
          </div>
        </div>
        <AnimatePresence initial={false}>
          {!isCollapsed && (
            <motion.div
              key="logs"
              initial={{ height: 0 }}
              animate={{ height: "auto" }}
              exit={{ height: 0 }}
              transition={{ duration: 0.15, ease: "easeInOut" }}
              className="overflow-hidden"
            >
              <div className="bg-gray-2 border-t-mauve-6 rounded-b-md border-t p-4">
                {!hasReceivedLogs && isBuilding ? (
                  <BuildCardSkeleton />
                ) : logs.length === 0 ? (
                  <pre className="text-mauve-11 max-h-125 overflow-y-auto whitespace-pre-wrap">
                    No logs available
                  </pre>
                ) : (
                  <pre
                    ref={logsScrollRef}
                    className="text-mauve-11 max-h-125 w-full overflow-y-auto font-mono text-sm break-all whitespace-pre-wrap"
                  >
                    {logs.map((line, i) => (
                      <span key={i} className="block">
                        {line}
                      </span>
                    ))}
                  </pre>
                )}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

function BuildCardSkeleton() {
  return (
    <div className="flex flex-col gap-1">
      <Skeleton className="h-5 w-96" />
      <Skeleton className="h-5 w-80" />
      <Skeleton className="h-5 w-52" />
      <Skeleton className="h-5 w-86" />
      <Skeleton className="h-5 w-24" />
      <Skeleton className="h-5 w-64" />
      <Skeleton className="h-5 w-72" />
    </div>
  );
}

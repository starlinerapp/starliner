import React from "react";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronRight } from "~/components/atoms/icons";
import { formatDistanceToNow } from "date-fns";
import { Spinner } from "~/components/atoms/spinner/Spinner";
import { Check, GitMerge, X } from "lucide-react";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";

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

export default function BuildCard({
  isCollapsed: collapsed = true,
  buildId,
  commitHash,
  source,
  serviceName,
  status,
  createdAt,
  args,
}: LogsCardProps) {
  const trpc = useTRPC();

  const [isCollapsed, setIsCollapsed] = React.useState(collapsed);
  const shouldPoll = status === "building" || status === "queued";

  const { data: logsData, isLoading } = useQuery(
    trpc.build.getBuildLogs.queryOptions(
      {
        id: buildId,
      },
      {
        enabled: !isCollapsed,
        refetchInterval: shouldPoll ? 1000 : false,
      },
    ),
  );

  return (
    <div className="shadow-xs">
      <div className="border-mauve-6 rounded-t-md border px-4 py-3 text-sm">
        <div className="flex gap-3">
          <div className="flex h-5 w-5 shrink-0 items-center justify-center">
            {(status === "queued" || status === "building") && (
              <Spinner className="stroke-violet-10 size-5" />
            )}
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
              <p>
                {source === "manual" ? "Manually triggered" : "Commit to main"}
              </p>
              {args && args.length > 0 && (
                <p className="mt-1">
                  Build Args:{" "}
                  {args.map((arg, i) => (
                    <span key={arg.name} className="text-mauve-11">
                      {arg.name}
                      {i < args.length - 1 ? ", " : ""}
                    </span>
                  ))}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>
      <div className="border-mauve-6 rounded-b-md border-x-1 border-b-1 text-sm">
        <div
          className="flex cursor-pointer items-center gap-3 px-4 py-3 text-sm"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          <motion.div
            animate={{ rotate: isCollapsed ? 0 : 90 }}
            transition={{ duration: 0.2, ease: "easeOut" }}
          >
            <ChevronRight className="w-4 stroke-2" />
          </motion.div>
          <p>Build Logs</p>
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
              <div className="bg-gray-2 border-t-mauve-6 border-t p-4">
                {isLoading || !logsData?.logs ? (
                  <BuildCardSkeleton />
                ) : (
                  <pre className="text-mauve-11 max-h-[500px] overflow-y-auto whitespace-pre-wrap">
                    {logsData?.logs || "No logs available"}
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

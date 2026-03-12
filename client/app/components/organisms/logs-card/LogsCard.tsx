import React from "react";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronRight } from "~/components/atoms/icons";
import { formatDistanceToNow } from "date-fns";
import { Spinner } from "~/components/atoms/spinner/Spinner";
import { Check, X } from "lucide-react";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

interface LogsCardProps {
  serviceName: string;
  createdAt: string;
  status: string;
}

export default function LogsCard({
  serviceName,
  status,
  createdAt,
}: LogsCardProps) {
  const [isCollapsed, setIsCollapsed] = React.useState(true);

  return (
    <div className="shadow-xs">
      <div className="border-mauve-6 rounded-t-md border-1 px-4 py-3 text-sm">
        <div className="flex items-center gap-3">
          {(status === "queued" || status === "building") && (
            <Spinner className="stroke-violet-10 size-5" />
          )}
          {status === "success" && (
            <div className="bg-grass-9 flex h-4.5 w-4.5 items-center justify-center rounded-full">
              <Check className="w-3.5 stroke-white stroke-2" />
            </div>
          )}
          {status === "error" && (
            <div className="bg-red-9 flex h-4.5 w-4.5 items-center justify-center rounded-full">
              <X className="w-3.5 stroke-white stroke-2" />
            </div>
          )}
          <div className="flex gap-2">
            <p className="font-bold">{serviceName}</p>
            <p>·</p>
            <p className="text-mauve-10">
              {formatDistanceToNow(createdAt)} ago
            </p>
          </div>
        </div>
      </div>
      <div className="border-mauve-6 rounded-b-md border-x-1 border-b-1 text-sm">
        <div
          className="flex cursor-pointer items-center gap-2 px-4 py-3 text-sm"
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
              animate={{ height: 200 }}
              exit={{ height: 0 }}
              transition={{ duration: 0.15, ease: "easeInOut" }}
              className="overflow-hidden"
            >
              <div className="bg-gray-2 border-t-mauve-6 h-[200px] border-t p-4">
                {true ? (
                  <LogsCardSkeleton />
                ) : (
                  <pre className="text-mauve-11 whitespace-pre-wrap">
                    Running build in Washington, D.C., USA (East) – iad1 Build
                    machine
                    {"\n"}
                    configuration: 2 cores, 8 GB Cloning
                    {"\n"}
                    github.com/leon-liang/team-server-down (Branch: develop,
                    Commit:
                    {"\n"}
                    9e18c70) Previous build caches not available Cloning
                    completed:
                    {"\n"}
                    1.396s Running "vercel build" Vercel CLI 48.1.6 Build
                    Completed in
                    {"\n"}
                    /vercel/output [222ms] Deploying outputs... Deployment
                    completed
                    {"\n"}
                    Creating build cache...
                    {"\n"}
                    Skipping cache upload because no files were prepared
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

function LogsCardSkeleton() {
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

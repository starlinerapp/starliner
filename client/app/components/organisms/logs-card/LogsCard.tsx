import React from "react";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronRight } from "~/components/atoms/icons";

interface LogsCardProps {
  serviceName: string;
}

export default function LogsCard({ serviceName }: LogsCardProps) {
  const [isCollapsed, setIsCollapsed] = React.useState(true);

  return (
    <div className="shadow-xs">
      <div className="border-mauve-6 rounded-t-md border-1 px-4 py-3 text-sm">
        <div className="flex gap-2">
          <p className="font-bold">{serviceName}</p>
          <p>·</p>
          <p className="text-mauve-10">20 min ago</p>
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
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

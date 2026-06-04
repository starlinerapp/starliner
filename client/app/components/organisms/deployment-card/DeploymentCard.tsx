import React, {
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from "react";
import { motion } from "framer-motion";
import { ChevronRight } from "~/components/atoms/icons";
import { formatDistanceToNow } from "date-fns";
import { Spinner } from "~/components/atoms/spinner/Spinner";
import { Check, GitMerge, X } from "lucide-react";
import { BuildLogs, BuildTab } from "./Build";
import {
  DeploymentLogs,
  DeploymentTab,
} from "~/components/organisms/deployment-card/Deployment";
import { cn } from "~/utils/cn";
import { scrollContainerToTop } from "./scroll";

interface LogsCardProps {
  isCollapsed?: boolean;
  buildId: number;
  deploymentId: number;
  commitHash: string;
  source: string;
  serviceName: string;
  createdAt: string;
  status: string;
  deploymentRolloutStatus: string;
}

const EXPAND_TRANSITION_MS = 200;

export default function DeploymentCard({
  isCollapsed: collapsed = true,
  buildId,
  deploymentId,
  commitHash,
  source,
  serviceName,
  status,
  deploymentRolloutStatus,
  createdAt,
}: LogsCardProps) {
  const isDeployOnly =
    source === "duplicate" ||
    (source === "manual" && status === "success" && !commitHash);
  const [isCollapsed, setIsCollapsed] = useState(collapsed);
  const [spacerReady, setSpacerReady] = useState(!collapsed);
  const [activePhase, setActivePhase] = useState<"build" | "deploy">(
    isDeployOnly ? "deploy" : "build",
  );
  const [hasBuildLogs, setHasBuildLogs] = useState(false);
  const [hasDeployLogs, setHasDeployLogs] = useState(false);
  const previousStatusRef = useRef(status);
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  const isBuilding =
    !isDeployOnly && (status === "queued" || status === "building");
  const isDeploying = isDeployOnly
    ? deploymentRolloutStatus === "pending"
    : status === "success" && deploymentRolloutStatus === "pending";
  const isComplete = isDeployOnly
    ? deploymentRolloutStatus === "success"
    : status === "success" && deploymentRolloutStatus === "success";
  const isFailed = isDeployOnly
    ? deploymentRolloutStatus === "failure"
    : status === "failure" || deploymentRolloutStatus === "failure";
  const showSpinner = isBuilding || isDeploying;

  const scrollToPhase = useCallback((behavior: ScrollBehavior = "smooth") => {
    const container = scrollContainerRef.current;
    if (!container) {
      return;
    }

    scrollContainerToTop(container, behavior);
  }, []);

  const expandCard = useCallback(() => {
    setIsCollapsed(false);
  }, []);

  const toggleCollapsed = useCallback(() => {
    setIsCollapsed((prev) => !prev);
  }, []);

  useEffect(() => {
    if (isCollapsed) {
      setSpacerReady(false);
      return;
    }

    const timeout = window.setTimeout(
      () => setSpacerReady(true),
      EXPAND_TRANSITION_MS,
    );
    return () => window.clearTimeout(timeout);
  }, [isCollapsed]);

  useLayoutEffect(() => {
    if (isCollapsed || !spacerReady) {
      return;
    }

    scrollToPhase("smooth");
  }, [activePhase, isCollapsed, spacerReady, scrollToPhase]);

  useEffect(() => {
    const previousStatus = previousStatusRef.current;
    previousStatusRef.current = status;

    if (
      status === "success" &&
      (previousStatus === "building" || previousStatus === "queued")
    ) {
      setActivePhase("deploy");
      expandCard();
    }
  }, [status, expandCard]);

  const selectBuild = () => {
    setActivePhase("build");
    if (isCollapsed) {
      expandCard();
    }
  };

  const selectDeploy = () => {
    setActivePhase("deploy");
    if (isCollapsed) {
      expandCard();
    }
  };

  return (
    <div className="shadow-xs">
      <div className="border-mauve-6 rounded-t-md border px-4 py-3 text-sm">
        <div className="flex gap-3">
          <div className="flex h-5 w-5 shrink-0 items-center justify-center">
            {showSpinner && <Spinner className="stroke-violet-10 size-5" />}
            {isComplete && (
              <div className="bg-grass-9 flex h-4.5 w-4.5 items-center justify-center rounded-full">
                <Check className="w-3.5 stroke-white stroke-2" />
              </div>
            )}
            {isFailed && (
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
                {source === "duplicate"
                  ? "On environment clone"
                  : source === "manual"
                    ? "Manually triggered"
                    : "On Push"}
              </p>
            </div>
          </div>
        </div>
      </div>
      <div className="border-mauve-6 rounded-b-md border-x border-b text-sm">
        <div
          className="flex cursor-pointer items-center gap-3 px-4 py-2 text-sm"
          onClick={toggleCollapsed}
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
            <BuildTab
              isActive={!isCollapsed && activePhase === "build"}
              hasLogs={hasBuildLogs || isDeployOnly}
              onSelect={selectBuild}
            />
            <div className="bg-mauve-8 h-px w-4" />
            <DeploymentTab
              isActive={!isCollapsed && activePhase === "deploy"}
              hasLogs={hasDeployLogs}
              onSelect={selectDeploy}
            />
          </div>
        </div>
        <div
          className={cn(
            "grid transition-[grid-template-rows] ease-in-out",
            isCollapsed ? "grid-rows-[0fr]" : "grid-rows-[1fr]",
          )}
          style={{ transitionDuration: `${EXPAND_TRANSITION_MS}ms` }}
        >
          <div className="min-h-0 overflow-hidden">
            <div
              ref={scrollContainerRef}
              className="bg-gray-2 border-t-mauve-6 overflow-anchor-none max-h-125 overflow-y-auto scroll-smooth rounded-b-md border-t p-4"
            >
              <div className={cn(activePhase === "deploy" && "hidden")}>
                {isDeployOnly ? (
                  <pre className="text-mauve-11 text-sm whitespace-pre-wrap">
                    Build step skipped
                  </pre>
                ) : (
                  <BuildLogs
                    buildId={buildId}
                    followScroll={activePhase === "build" && isBuilding}
                    onHasLogsChange={setHasBuildLogs}
                  />
                )}
              </div>
              <div
                className={cn(
                  !isBuilding &&
                    activePhase === "build" &&
                    "border-mauve-6 mt-4 border-t pt-4",
                )}
              >
                <DeploymentLogs
                  deploymentId={deploymentId}
                  buildStatus={status}
                  followScroll={activePhase === "deploy" && isDeploying}
                  onHasLogsChange={setHasDeployLogs}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

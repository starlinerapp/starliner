import { formatDistanceToNow } from "date-fns";
import { motion } from "framer-motion";
import { Check, GitMerge, X } from "lucide-react";
import {
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from "react";
import { ChevronRight } from "~/components/atoms/icons";
import { Spinner } from "~/components/atoms/spinner/Spinner";
import {
  DeploymentLogs,
  DeploymentTab,
} from "~/components/organisms/deployment-card/Deployment";
import { cn } from "~/utils/cn";
import { BuildLogs, BuildTab } from "./Build";
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
      <div className="rounded-t-md border border-mauve-6 px-4 py-3 text-sm">
        <div className="flex gap-3">
          <div className="flex h-5 w-5 shrink-0 items-center justify-center">
            {showSpinner && <Spinner className="size-5 stroke-violet-10" />}
            {isComplete && (
              <div className="flex h-4.5 w-4.5 items-center justify-center rounded-full bg-grass-9">
                <Check className="w-3.5 stroke-2 stroke-white" />
              </div>
            )}
            {isFailed && (
              <div className="flex h-4.5 w-4.5 items-center justify-center rounded-full bg-red-9">
                <X className="w-3.5 stroke-2 stroke-white" />
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
                <p className="flex items-center gap-1 rounded-md border border-mauve-6 bg-gray-3 px-1.5 text-mauve-10">
                  <GitMerge size={16} />
                  {commitHash.slice(0, 7)}
                </p>
              )}
            </div>

            <div className="mt-0.5 text-mauve-10 text-sm">
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
      <div className="rounded-b-md border-mauve-6 border-x border-b text-sm">
        <div className="relative flex items-center gap-3 px-4 py-2 text-sm">
          <button
            type="button"
            aria-label="Toggle logs"
            aria-expanded={!isCollapsed}
            onClick={toggleCollapsed}
            className="absolute z-10 inset-0 cursor-pointer"
          />
          <motion.div
            className="relative"
            animate={{ rotate: isCollapsed ? 0 : 90 }}
            transition={{ duration: 0.2, ease: "easeOut" }}
          >
            <ChevronRight className="w-4 stroke-2" />
          </motion.div>
          <div className="relative flex items-center">
            <BuildTab
              isActive={!isCollapsed && activePhase === "build"}
              hasLogs={!isDeployOnly && hasBuildLogs}
              onSelect={selectBuild}
            />
            <div className="h-px w-4 bg-mauve-8" />
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
              className="overflow-anchor-none max-h-125 overflow-y-auto scroll-smooth rounded-b-md border-t border-t-mauve-6 bg-gray-2 p-4"
            >
              <div className={cn(activePhase === "deploy" && "hidden")}>
                {isDeployOnly ? (
                  <pre className="whitespace-pre-wrap text-mauve-11 text-sm">
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
                    "mt-4 border-mauve-6 border-t pt-4",
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

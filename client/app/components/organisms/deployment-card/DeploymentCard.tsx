import React, { useCallback, useEffect, useRef, useState } from "react";
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
import {
  scrollContainerToSectionBottom,
  scrollContainerToSectionStart,
  scrollContainerToTop,
} from "./scroll";

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
  autoSwitchToDeployAfterBuildLogs?: boolean;
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
  autoSwitchToDeployAfterBuildLogs = false,
}: LogsCardProps) {
  const isDeployOnly = source === "duplicate";
  const [isCollapsed, setIsCollapsed] = useState(collapsed);
  const [spacerReady, setSpacerReady] = useState(!collapsed);
  const [activePhase, setActivePhase] = useState<"build" | "deploy">(
    isDeployOnly ? "deploy" : "build",
  );
  const [hasBuildLogs, setHasBuildLogs] = useState(false);
  const [hasDeployLogs, setHasDeployLogs] = useState(false);
  const [bottomSpacerHeight, setBottomSpacerHeight] = useState(0);
  const previousStatusRef = useRef(status);
  const wasCollapsedRef = useRef(isCollapsed);
  const activePhaseRef = useRef(activePhase);
  const hasAutoSwitchedToDeployRef = useRef(false);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const buildSectionRef = useRef<HTMLDivElement>(null);
  const deploySectionRef = useRef<HTMLDivElement>(null);

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

  activePhaseRef.current = activePhase;

  const updateBottomSpacer = useCallback(() => {
    const container = scrollContainerRef.current;
    const deploy = deploySectionRef.current;
    if (!container || !deploy) {
      return;
    }

    const containerHeight = container.clientHeight;
    const existingSpacer = deploy.querySelector<HTMLElement>(
      "[data-deploy-scroll-spacer]",
    );
    const existingSpacerHeight = existingSpacer?.offsetHeight ?? 0;
    const deployContentHeight = deploy.offsetHeight - existingSpacerHeight;

    setBottomSpacerHeight(Math.max(0, containerHeight - deployContentHeight));
  }, []);

  const scrollToPhase = useCallback(
    (phase: "build" | "deploy", behavior: ScrollBehavior = "smooth") => {
      const container = scrollContainerRef.current;
      const deploySection = deploySectionRef.current;
      if (!container) {
        return;
      }

      if (phase === "build") {
        scrollContainerToTop(container, behavior);
        return;
      }

      if (deploySection) {
        scrollContainerToSectionStart(container, deploySection, behavior);
      }
    },
    [],
  );

  const expandCard = useCallback(() => {
    setIsCollapsed(false);
  }, []);

  const toggleCollapsed = useCallback(() => {
    setIsCollapsed((prev) => !prev);
  }, []);

  useEffect(() => {
    if (isCollapsed) {
      setSpacerReady(false);
      setBottomSpacerHeight(0);
      return;
    }

    const timeout = window.setTimeout(
      () => setSpacerReady(true),
      EXPAND_TRANSITION_MS,
    );
    return () => window.clearTimeout(timeout);
  }, [isCollapsed]);

  useEffect(() => {
    if (isCollapsed || !spacerReady) {
      return;
    }

    updateBottomSpacer();

    const container = scrollContainerRef.current;
    const build = buildSectionRef.current;
    const deploy = deploySectionRef.current;
    if (!container) {
      return;
    }

    const observer = new ResizeObserver(() => {
      updateBottomSpacer();
    });

    observer.observe(container);
    if (build) {
      observer.observe(build);
    }
    if (deploy) {
      observer.observe(deploy);
    }

    return () => observer.disconnect();
  }, [isCollapsed, spacerReady, updateBottomSpacer, buildId, deploymentId]);

  useEffect(() => {
    const wasCollapsed = wasCollapsedRef.current;
    wasCollapsedRef.current = isCollapsed;
    if (!wasCollapsed || isCollapsed) {
      return;
    }

    requestAnimationFrame(() => {
      const container = scrollContainerRef.current;
      if (activePhaseRef.current === "build" && container) {
        scrollContainerToTop(container, "instant");
      }
    });
  }, [isCollapsed]);

  useEffect(() => {
    if (!spacerReady || activePhase !== "deploy" || isCollapsed) {
      return;
    }
    requestAnimationFrame(() => {
      scrollToPhase("deploy", "smooth");
    });
  }, [spacerReady, activePhase, isCollapsed, scrollToPhase]);

  useEffect(() => {
    hasAutoSwitchedToDeployRef.current = false;
  }, [buildId]);

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

  useEffect(() => {
    if (
      !autoSwitchToDeployAfterBuildLogs ||
      hasAutoSwitchedToDeployRef.current
    ) {
      return;
    }
    if (status !== "success" || !hasBuildLogs) {
      return;
    }

    hasAutoSwitchedToDeployRef.current = true;
    setActivePhase("deploy");
    expandCard();
  }, [autoSwitchToDeployAfterBuildLogs, status, hasBuildLogs, expandCard]);

  const selectBuild = () => {
    setActivePhase("build");
    if (isCollapsed) {
      expandCard();
    }
    requestAnimationFrame(() => {
      scrollToPhase("build", "smooth");
    });
  };

  const selectDeploy = () => {
    setActivePhase("deploy");
    if (isCollapsed) {
      expandCard();
      return;
    }
    requestAnimationFrame(() => {
      scrollToPhase("deploy", "smooth");
    });
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
              <div ref={buildSectionRef}>
                {isDeployOnly ? (
                  <pre className="text-mauve-11 text-sm whitespace-pre-wrap">
                    Build step skipped
                  </pre>
                ) : (
                  <BuildLogs
                    buildId={buildId}
                    followScroll={activePhase === "build" && isBuilding}
                    scrollContainerRef={scrollContainerRef}
                    sectionRef={buildSectionRef}
                    onHasLogsChange={setHasBuildLogs}
                  />
                )}
              </div>
              <div
                ref={deploySectionRef}
                className={cn(
                  !isBuilding && "border-mauve-6 mt-4 border-t pt-4",
                )}
              >
                <DeploymentLogs
                  deploymentId={deploymentId}
                  buildStatus={status}
                  followScroll={activePhase === "deploy" && isDeploying}
                  scrollContainerRef={scrollContainerRef}
                  sectionRef={deploySectionRef}
                  onHasLogsChange={setHasDeployLogs}
                />
                {spacerReady && bottomSpacerHeight > 0 && (
                  <div
                    aria-hidden
                    data-deploy-scroll-spacer
                    className="shrink-0"
                    style={{ height: bottomSpacerHeight }}
                  />
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

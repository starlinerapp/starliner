import React, { useEffect, useState } from "react";
import { Play } from "lucide-react";
import { useSubscription } from "@trpc/tanstack-react-query";
import { useTRPC } from "~/utils/trpc/react";
import { cn } from "~/utils/cn";
import { scrollContainerToSectionBottom } from "./scroll";

interface DeploymentTabProps {
  isActive: boolean;
  hasLogs: boolean;
  onSelect: () => void;
}

export function DeploymentTab({
  isActive,
  hasLogs,
  onSelect,
}: DeploymentTabProps) {
  return (
    <div className="relative">
      <div className="bg-mauve-8 absolute top-1/2 -left-1 h-2 w-2 -translate-y-1/2 rounded-full" />
      <button
        type="button"
        onClick={onSelect}
        className={cn(
          "hover:bg-mauve-2 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border bg-white px-4 py-0.5",
          !hasLogs && "border-mauve-6 text-mauve-8",
          hasLogs && isActive && "border-violet-9 bg-violet-3 text-violet-9",
          hasLogs && !isActive && "border-mauve-9 text-mauve-9",
        )}
      >
        <div
          className={cn(
            "flex rounded-full border-[1.5px] p-0.5",
            !hasLogs && "border-mauve-8",
            hasLogs && isActive && "border-violet-9",
            hasLogs && !isActive && "border-mauve-9",
          )}
        >
          <Play
            className={cn(
              "h-2 w-2",
              !hasLogs && "fill-mauve-8",
              hasLogs && isActive && "fill-violet-9 stroke-violet-9",
              hasLogs && !isActive && "fill-mauve-9 stroke-mauve-9",
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
  followScroll?: boolean;
  scrollContainerRef?: React.RefObject<HTMLDivElement | null>;
  sectionRef?: React.RefObject<HTMLDivElement | null>;
  onHasLogsChange?: (hasLogs: boolean) => void;
}

export function DeploymentLogs({
  deploymentId,
  buildStatus,
  followScroll = false,
  scrollContainerRef,
  sectionRef,
  onHasLogsChange,
}: DeploymentLogsProps) {
  const trpc = useTRPC();
  const [lines, setLines] = useState<string[]>([]);

  useEffect(() => {
    setLines([]);
  }, [deploymentId]);

  const buildComplete = buildStatus === "success";
  const buildFailed = buildStatus === "failure";

  useSubscription(
    trpc.deployment.streamDeploymentStatusLogs.subscriptionOptions(
      { deploymentId },
      {
        enabled: buildComplete,
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

  useEffect(() => {
    onHasLogsChange?.(lines.length > 0);
  }, [lines, onHasLogsChange]);

  useEffect(() => {
    if (!followScroll || !scrollContainerRef?.current || !sectionRef?.current) {
      return;
    }

    scrollContainerToSectionBottom(
      scrollContainerRef.current,
      sectionRef.current,
      "smooth",
      "[data-deploy-scroll-spacer]",
    );
  }, [lines, followScroll, scrollContainerRef, sectionRef]);

  if (buildFailed) {
    return (
      <pre className="text-mauve-11 text-sm whitespace-pre-wrap">
        Build failed — deployment was not triggered.
      </pre>
    );
  }

  return (
    <div className="text-mauve-11 font-mono text-sm break-all whitespace-pre-wrap">
      {lines.map((line, i) =>
        line === "" ? (
          <span key={i} className="block h-4" aria-hidden />
        ) : (
          <span key={i} className="block">
            {line}
          </span>
        ),
      )}
    </div>
  );
}

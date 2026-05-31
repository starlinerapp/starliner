import React, { useEffect, useState } from "react";
import { Hammer } from "lucide-react";
import { useSubscription } from "@trpc/tanstack-react-query";
import { useTRPC } from "~/utils/trpc/react";
import { cn } from "~/utils/cn";

interface BuildTabProps {
  isActive: boolean;
  hasLogs: boolean;
  onSelect: () => void;
}

export function BuildTab({ isActive, hasLogs, onSelect }: BuildTabProps) {
  return (
    <div className="relative">
      <div className="bg-mauve-8 absolute top-1/2 -right-1 h-2 w-2 -translate-y-1/2 rounded-full" />
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
          <Hammer
            className={cn(
              "h-2 w-2",
              !hasLogs && "fill-mauve-8",
              hasLogs && isActive && "fill-violet-9 stroke-violet-9",
              hasLogs && !isActive && "fill-mauve-9 stroke-mauve-9",
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
  followScroll?: boolean;
  scrollContainerRef?: React.RefObject<HTMLDivElement | null>;
  sectionRef?: React.RefObject<HTMLDivElement | null>;
  onHasLogsChange?: (hasLogs: boolean) => void;
}

export function BuildLogs({
  buildId,
  followScroll = false,
  scrollContainerRef,
  sectionRef,
  onHasLogsChange,
}: BuildLogsProps) {
  const trpc = useTRPC();
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    setLogs([]);
  }, [buildId]);

  useSubscription(
    trpc.build.streamBuildLogs.subscriptionOptions(
      { buildId },
      {
        onData: (chunk) => {
          setLogs((prev) => [...prev, chunk]);
        },
      },
    ),
  );

  useEffect(() => {
    onHasLogsChange?.(logs.length > 0);
  }, [logs, onHasLogsChange]);

  useEffect(() => {
    if (!followScroll || !scrollContainerRef?.current) {
      return;
    }
    const container = scrollContainerRef.current;
    const section = sectionRef?.current;
    const targetScroll = section
      ? Math.max(
          0,
          section.offsetTop + section.offsetHeight - container.clientHeight,
        )
      : container.scrollHeight - container.clientHeight;

    container.scrollTo({ top: targetScroll, behavior: "smooth" });
  }, [logs, followScroll, scrollContainerRef, sectionRef]);

  return (
    <div className="text-mauve-11 font-mono text-sm break-all whitespace-pre-wrap">
      {logs.map((line, i) => (
        <span key={i} className="block">
          {line}
        </span>
      ))}
    </div>
  );
}

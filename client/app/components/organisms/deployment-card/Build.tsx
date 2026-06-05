import { useSubscription } from "@trpc/tanstack-react-query";
import { Hammer } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

interface BuildTabProps {
  isActive: boolean;
  hasLogs: boolean;
  onSelect: () => void;
}

export function BuildTab({ isActive, hasLogs, onSelect }: BuildTabProps) {
  return (
    <div className="relative">
      <div className="absolute top-1/2 -right-1 h-2 w-2 -translate-y-1/2 rounded-full bg-mauve-8" />
      <button
        type="button"
        onClick={onSelect}
        className={cn(
          "relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border bg-white px-4 py-0.5 hover:bg-mauve-2",
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
  onHasLogsChange?: (hasLogs: boolean) => void;
}

export function BuildLogs({
  buildId,
  followScroll = false,
  onHasLogsChange,
}: BuildLogsProps) {
  const trpc = useTRPC();
  const [logs, setLogs] = useState<string[]>([]);
  const tailRef = useRef<HTMLSpanElement>(null);

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
    if (!followScroll || !tailRef.current) {
      return;
    }

    tailRef.current.scrollIntoView({ behavior: "smooth", block: "end" });
  }, [logs, followScroll]);

  return (
    <div className="whitespace-pre-wrap break-all font-mono text-mauve-11 text-sm">
      {logs.map((line, i) => (
        <span key={i} className="block">
          {line}
        </span>
      ))}
      <span ref={tailRef} className="block h-px" aria-hidden />
    </div>
  );
}

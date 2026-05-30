import React, { useEffect, useRef, useState } from "react";
import { Hammer } from "lucide-react";
import { useSubscription } from "@trpc/tanstack-react-query";
import { useTRPC } from "~/utils/trpc/react";

interface BuildTabProps {
  isActive: boolean;
  onSelect: () => void;
}

export function BuildTab({ isActive, onSelect }: BuildTabProps) {
  return (
    <div className="relative">
      <div className="bg-mauve-8 absolute top-1/2 -right-1 h-2 w-2 -translate-y-1/2 rounded-full" />
      <button
        type="button"
        onClick={onSelect}
        className={
          isActive
            ? "border-violet-9 bg-violet-3 hover:bg-mauve-2 text-violet-9 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border px-4 py-0.5"
            : "border-mauve-6 hover:bg-mauve-2 text-mauve-9 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border bg-white px-4 py-0.5"
        }
      >
        <div
          className={
            isActive
              ? "border-violet-9 flex rounded-full border-[1.5px] p-0.5"
              : "border-mauve-9 flex rounded-full border-[1.5px] p-0.5"
          }
        >
          <Hammer
            className={
              isActive
                ? "fill-violet-9 stroke-violet-9 h-2 w-2"
                : "fill-mauve-9 stroke-mauve-9 h-2 w-2"
            }
          />
        </div>
        Build
      </button>
    </div>
  );
}

interface BuildLogsProps {
  buildId: number;
  status: string;
  loadingFallback: React.ReactNode;
}

export function BuildLogs({
  buildId,
  status,
  loadingFallback,
}: BuildLogsProps) {
  const trpc = useTRPC();
  const [logs, setLogs] = useState<string[]>([]);
  const [hasReceivedLogs, setHasReceivedLogs] = useState(false);
  const logsScrollRef = useRef<HTMLPreElement>(null);
  const hasLoadedInitial = useRef(false);

  useEffect(() => {
    setLogs([]);
    setHasReceivedLogs(false);
    hasLoadedInitial.current = false;
  }, [buildId]);

  useSubscription(
    trpc.build.streamBuildLogs.subscriptionOptions(
      { buildId },
      {
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

  if (!hasReceivedLogs && isBuilding) {
    return loadingFallback;
  }

  if (logs.length === 0) {
    return (
      <pre className="text-mauve-11 max-h-125 overflow-y-auto whitespace-pre-wrap">
        No logs available
      </pre>
    );
  }

  return (
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
  );
}

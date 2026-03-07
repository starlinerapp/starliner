import React, { useEffect, useRef, useState } from "react";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/client/generated";
import { useTRPC } from "~/utils/trpc/react";
import { useSubscription } from "@trpc/tanstack-react-query";

type Deployment =
  | ResponseGitDeployment
  | ResponseImageDeployment
  | ResponseIngressDeployment
  | ResponseDatabaseDeployment;

interface BottomBarProps {
  deployment?: Deployment;
}

export default function BottomBar({ deployment }: BottomBarProps) {
  const trpc = useTRPC();

  const containerRef = useRef<HTMLDivElement>(null);
  const titleRef = useRef<HTMLSpanElement>(null);
  const logsEndRef = useRef<HTMLDivElement>(null);
  const hasLoadedInitial = useRef(false);

  const [underline, setUnderline] = useState({ left: 0, width: 0 });
  const [logs, setLogs] = useState<string[]>([]);

  useSubscription(
    trpc.deployment.streamDeploymentLogs.subscriptionOptions(
      { deploymentId: Number(deployment?.id) },
      {
        enabled: !!deployment?.id,
        onData: (chunk) => setLogs((prev) => [...prev, chunk]),
      },
    ),
  );

  useEffect(() => {
    if (deployment) {
      hasLoadedInitial.current = false;
      setLogs([]);
    }
  }, [deployment]);

  useEffect(() => {
    if (!containerRef.current || !titleRef.current) return;
    const rect = titleRef.current.getBoundingClientRect();
    const containerRect = containerRef.current.getBoundingClientRect();
    setUnderline({
      left: rect.left - containerRect.left,
      width: rect.width,
    });
  }, [deployment]);

  useEffect(() => {
    if (!hasLoadedInitial.current) {
      if (logs.length > 0) {
        hasLoadedInitial.current = true;
        requestAnimationFrame(() => {
          requestAnimationFrame(() => {
            logsEndRef.current?.scrollIntoView({ behavior: "instant" });
          });
        });
      }
      return;
    }
    logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <div className="flex h-full flex-col">
      <div className="bg-violet-1 flex-shrink-0">
        <div
          ref={containerRef}
          className="border-mauve-6 text-mauve-11 relative flex w-full gap-4 border-b px-2 pt-2 pb-1 text-sm"
        >
          <div className="relative z-10 px-2 py-1.5">
            <span
              ref={titleRef}
              className="text-violet-11 truncate pb-2 font-semibold"
            >
              Logs
            </span>
          </div>
          <div
            className="bg-violet-11 absolute bottom-0 h-[3px] rounded-md"
            style={{ left: underline.left, width: underline.width }}
          />
        </div>
      </div>

      <div className="min-h-0 flex-1 overflow-y-scroll p-4">
        {!deployment ? (
          <p className="text-mauve-11">
            No deployment selected. Select one to view logs.
          </p>
        ) : (
          <pre className="text-mauve-11 w-full font-mono text-sm break-all whitespace-pre-wrap">
            {logs.map((line, i) => (
              <span key={i} className="block">
                {line}
              </span>
            ))}
            <div ref={logsEndRef} />
          </pre>
        )}
      </div>
    </div>
  );
}

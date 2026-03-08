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

interface LogsProps {
  deployment: Deployment | undefined;
}

export default function Logs({ deployment }: LogsProps) {
  const trpc = useTRPC();

  const hasLoadedInitial = useRef(false);

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

  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (deployment) {
      hasLoadedInitial.current = false;
      setLogs([]);
    }
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
    <>
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
    </>
  );
}

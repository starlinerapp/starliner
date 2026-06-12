import { useSubscription } from "@trpc/tanstack-react-query";
import { useEffect, useState } from "react";
import LogsViewer from "~/components/molecules/logs-viewer/LogsViewer";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/clients/server/generated";
import { useTRPC } from "~/utils/trpc/react";

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
      setLogs([]);
    }
  }, [deployment]);

  return (
    <>
      {!deployment ? (
        <p className="text-mauve-11">
          No deployment selected. Select one to view logs.
        </p>
      ) : (
        <LogsViewer logs={logs} resetKey={deployment.id} />
      )}
    </>
  );
}

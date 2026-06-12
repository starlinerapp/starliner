import { useSubscription } from "@trpc/tanstack-react-query";
import { useEffect, useState } from "react";
import LogsViewer from "~/components/molecules/logs-viewer/LogsViewer";
import { useTRPC } from "~/utils/trpc/react";

interface LogsProps {
  clusterId: number | undefined;
}

export default function Logs({ clusterId }: LogsProps) {
  const trpc = useTRPC();

  const [logs, setLogs] = useState<string[]>([]);

  useSubscription(
    trpc.cluster.streamProvisioningLogs.subscriptionOptions(
      { clusterId: Number(clusterId) },
      {
        enabled: !!clusterId,
        onData: (chunk) => setLogs((prev) => [...prev, chunk]),
      },
    ),
  );

  useEffect(() => {
    if (clusterId) {
      setLogs([]);
    }
  }, [clusterId]);

  return (
    <>
      {!clusterId ? (
        <p className="text-mauve-11">No cluster selected.</p>
      ) : (
        <LogsViewer logs={logs} resetKey={clusterId} />
      )}
    </>
  );
}

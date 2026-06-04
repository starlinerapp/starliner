import { useSubscription } from "@trpc/tanstack-react-query";
import { useEffect, useRef, useState } from "react";
import { useTRPC } from "~/utils/trpc/react";

interface LogsProps {
  clusterId: number | undefined;
}

export default function Logs({ clusterId }: LogsProps) {
  const trpc = useTRPC();

  const hasLoadedInitial = useRef(false);

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

  const logsScrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (clusterId) {
      hasLoadedInitial.current = false;
      setLogs([]);
    }
  }, [clusterId]);

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

  return (
    <>
      {!clusterId ? (
        <p className="text-mauve-11">No cluster selected.</p>
      ) : (
        <div
          ref={logsScrollRef}
          className="h-full min-h-0 w-full overflow-y-auto"
        >
          <pre className="w-full whitespace-pre-wrap break-all font-mono text-mauve-11 text-sm">
            {logs.map((line, i) => (
              <span key={i} className="block">
                {line}
              </span>
            ))}
          </pre>
        </div>
      )}
    </>
  );
}

import { useSubscription } from "@trpc/tanstack-react-query";
import { useMemo } from "react";
import { useToast } from "~/components/atoms/toast/ToastProvider";
import { useTRPC } from "~/utils/trpc/react";

export function useEnvironmentDeploymentNotifications(environmentId?: number) {
  const trpc = useTRPC();
  const toast = useToast();

  const subscriptionOptions = useMemo(
    () =>
      trpc.environment.streamDeploymentNotifications.subscriptionOptions(
        { id: environmentId ?? 0 },
        {
          enabled: environmentId != null,
          onData: (data: { status: string; message: string }) => {
            if (data.status === "success") {
              toast.success(data.message);
            } else if (data.status === "failed") {
              toast.error(data.message);
            }
          },
        },
      ),
    [trpc, environmentId, toast],
  );

  useSubscription(subscriptionOptions);
}

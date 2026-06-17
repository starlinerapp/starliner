import { useSubscription } from "@trpc/tanstack-react-query";
import { useMemo } from "react";
import { useToast } from "~/components/atoms/toast/ToastProvider";
import { useTRPC } from "~/utils/trpc/react";

export function useGlobalNotifications(organizationId: number) {
  const trpc = useTRPC();
  const toast = useToast();

  const subscriptionOptions = useMemo(
    () =>
      trpc.notifications.streamGlobalNotifications.subscriptionOptions(
        { organizationId },
        {
          onData: (data: {
            organizationId: number;
            status: string;
            message: string;
          }) => {
            if (data.status === "success") {
              toast.success(data.message);
            } else if (data.status === "failed") {
              toast.error(data.message);
            }
          },
        },
      ),
    [trpc, organizationId, toast],
  );

  useSubscription(subscriptionOptions);
}

import { useSubscription } from "@trpc/tanstack-react-query";
import { useMemo, useRef } from "react";
import {
  ErrorToast,
  type ToastHandle,
} from "~/components/atoms/toast/ErrorToast";
import { SuccessToast } from "~/components/atoms/toast/SuccessToast";
import { useTRPC } from "~/utils/trpc/react";

interface GlobalNotificationListenerProps {
  organizationId: number;
}

export default function GlobalNotificationListener({
  organizationId,
}: GlobalNotificationListenerProps) {
  const trpc = useTRPC();
  const successRef = useRef<ToastHandle>(null);
  const errorRef = useRef<ToastHandle>(null);

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
              successRef.current?.publish(data.message);
            } else if (data.status === "failed") {
              errorRef.current?.publish(data.message);
            }
          },
        },
      ),
    [organizationId],
  );

  useSubscription(subscriptionOptions);

  return (
    <>
      <SuccessToast ref={successRef} title="Success" />
      <ErrorToast ref={errorRef} title="Oops something went wrong" />
    </>
  );
}

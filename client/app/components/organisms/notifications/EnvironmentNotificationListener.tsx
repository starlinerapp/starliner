import React, { useMemo, useRef } from "react";
import { useSubscription } from "@trpc/tanstack-react-query";
import { useTRPC } from "~/utils/trpc/react";
import { SuccessToast } from "~/components/atoms/toast/SuccessToast";
import {
  ErrorToast,
  type ToastHandle,
} from "~/components/atoms/toast/ErrorToast";

interface EnvironmentNotificationListenerProps {
  environmentId: number;
}

export default function EnvironmentNotificationListener({
  environmentId,
}: EnvironmentNotificationListenerProps) {
  const trpc = useTRPC();
  const successRef = useRef<ToastHandle>(null);
  const errorRef = useRef<ToastHandle>(null);

  const subscriptionOptions = useMemo(
    () =>
      trpc.environment.streamDeploymentNotifications.subscriptionOptions(
        { id: environmentId },
        {
          onData: (data: { status: string; message: string }) => {
            if (data.status === "success") {
              successRef.current?.publish(data.message);
            } else if (data.status === "failed") {
              errorRef.current?.publish(data.message);
            }
          },
        },
      ),
    [environmentId],
  );

  useSubscription(subscriptionOptions);

  return (
    <>
      <SuccessToast ref={successRef} title="Success" />
      <ErrorToast ref={errorRef} title="Oops something went wrong" />
    </>
  );
}

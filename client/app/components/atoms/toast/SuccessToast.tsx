import * as React from "react";
import * as ToastPrimitive from "@radix-ui/react-toast";
import { cn } from "~/utils/cn";
import type { ToastHandle } from "~/components/atoms/toast/ErrorToast";
import { CheckCircle } from "~/components/atoms/icons";

interface SuccessToastProps {
  title?: string;
  duration?: number;
}

interface ToastItem {
  id: number;
  description: string;
}

export const SuccessToast = React.forwardRef<ToastHandle, SuccessToastProps>(
  ({ title = "Success!", duration = 4000 }, forwardedRef) => {
    const [toasts, setToasts] = React.useState<ToastItem[]>([]);
    const counterRef = React.useRef(0);

    React.useImperativeHandle(forwardedRef, () => ({
      publish: (message?: string) => {
        const id = ++counterRef.current;
        setToasts((prev) => [...prev, { id, description: message ?? "" }]);
      },
    }));

    const remove = (id: number) =>
      setToasts((prev) => prev.filter((t) => t.id !== id));

    return (
      <>
        {toasts.map((toast) => (
          <SingleSuccessToast
            key={toast.id}
            title={title}
            description={toast.description}
            duration={duration}
            onClose={() => remove(toast.id)}
          />
        ))}
      </>
    );
  },
);

SuccessToast.displayName = "SuccessToast";

function SingleSuccessToast({
  title,
  description,
  duration,
  onClose,
}: {
  title: string;
  description: string;
  duration: number;
  onClose: () => void;
}) {
  return (
    <ToastPrimitive.Root
      duration={duration}
      onOpenChange={(o) => {
        if (!o) onClose();
      }}
      className={cn(
        "border-grass-9 relative flex flex-row items-center gap-3 overflow-hidden rounded-xl border px-4 py-3 shadow-md",
        "data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:slide-in-from-bottom-4",
        "data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:slide-out-to-right-full",
        "data-[swipe=move]:translate-x-(--radix-toast-swipe-move-x)",
        "data-[swipe=cancel]:translate-x-0 data-[swipe=cancel]:transition-transform",
        "data-[swipe=end]:animate-out data-[swipe=end]:slide-out-to-right-full",
      )}
    >
      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg">
        <CheckCircle className="text-grass-10" />
      </div>
      <div className="flex flex-col gap-0.5">
        <ToastPrimitive.Title className="text-mauve-12 text-sm">
          {title}
        </ToastPrimitive.Title>
        <ToastPrimitive.Description className="text-mauve-10 text-sm">
          {description}
        </ToastPrimitive.Description>
      </div>
      <ToastPrimitive.Close className="sr-only" />
      <div className="absolute right-0 bottom-0 left-0 h-1">
        <div
          className="bg-grass-10 animate-shrink h-full origin-left"
          style={{ animation: `shrink ${duration}ms linear forwards` }}
        />
      </div>
    </ToastPrimitive.Root>
  );
}

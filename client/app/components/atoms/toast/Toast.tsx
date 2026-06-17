import * as ToastPrimitive from "@radix-ui/react-toast";
import { cva } from "class-variance-authority";
import { CheckCircle, ExclamationTriangle } from "~/components/atoms/icons";

export type ToastIntent = "success" | "error";

const toastVariants = cva(
  [
    "relative flex flex-row items-center gap-3 overflow-hidden rounded-xl border px-4 py-3 shadow-md",
    "data-[state=open]:fade-in-0 data-[state=open]:slide-in-from-bottom-4 data-[state=open]:animate-in",
    "data-[state=closed]:fade-out-0 data-[state=closed]:slide-out-to-right-full data-[state=closed]:animate-out",
    "data-[swipe=move]:translate-x-(--radix-toast-swipe-move-x)",
    "data-[swipe=cancel]:translate-x-0 data-[swipe=cancel]:transition-transform",
    "data-[swipe=end]:slide-out-to-right-full data-[swipe=end]:animate-out",
  ],
  {
    variants: {
      intent: {
        success: "border-grass-9",
        error: "border-red-9",
      },
    },
    defaultVariants: { intent: "success" },
  },
);

const iconVariants = cva("", {
  variants: {
    intent: {
      success: "text-grass-10",
      error: "text-red-9",
    },
  },
  defaultVariants: { intent: "success" },
});

const progressVariants = cva("h-full origin-left animate-shrink", {
  variants: {
    intent: {
      success: "bg-grass-10",
      error: "bg-red-9",
    },
  },
  defaultVariants: { intent: "success" },
});

const iconByIntent = {
  success: CheckCircle,
  error: ExclamationTriangle,
} as const;

interface ToastProps {
  intent: ToastIntent;
  title: string;
  description: string;
  duration: number;
  onClose: () => void;
}

export function Toast({
  intent,
  title,
  description,
  duration,
  onClose,
}: ToastProps) {
  const Icon = iconByIntent[intent];

  return (
    <ToastPrimitive.Root
      duration={duration}
      onOpenChange={(o) => {
        if (!o) onClose();
      }}
      className={toastVariants({ intent })}
    >
      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg">
        <Icon className={iconVariants({ intent })} />
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
          className={progressVariants({ intent })}
          style={{ animation: `shrink ${duration}ms linear forwards` }}
        />
      </div>
    </ToastPrimitive.Root>
  );
}

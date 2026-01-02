import React from "react";
import { cn } from "~/utils/cn";

interface LiveIndicatorProps {
  type: "warning" | "success" | "error";
}

export default function LiveIndicator({ type }: LiveIndicatorProps) {
  return (
    <span className="relative flex size-3">
      <span
        className={cn(
          "absolute inline-flex h-full w-full animate-ping rounded-full opacity-75",
          type === "warning" && "bg-orange-500",
          type === "success" && "bg-green-500",
          type === "error" && "bg-red-500",
        )}
      ></span>
      <span
        className={cn(
          "relative inline-flex size-3 rounded-full",
          type === "warning" && "bg-orange-400",
          type === "success" && "bg-green-400",
          type === "error" && "bg-red-400",
        )}
      ></span>
    </span>
  );
}

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
          type === "warning" && "bg-amber-10",
          type === "success" && "bg-grass-10",
          type === "error" && "bg-red-10",
        )}
      ></span>
      <span
        className={cn(
          "relative inline-flex size-3 rounded-full",
          type === "warning" && "bg-amber-9",
          type === "success" && "bg-grass-9",
          type === "error" && "bg-red-9",
        )}
      ></span>
    </span>
  );
}

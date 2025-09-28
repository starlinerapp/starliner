import React from "react";
import { cn } from "~/utils/cn";

export default function Skeleton({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="skeleton"
      className={cn("bg-gray-5 animate-pulse rounded-md", className)}
      {...props}
    />
  );
}
export { Skeleton };

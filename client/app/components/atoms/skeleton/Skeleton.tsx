import type React from "react";
import { cn } from "~/utils/cn";

export default function Skeleton({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="skeleton"
      className={cn("animate-pulse rounded-md bg-gray-5", className)}
      {...props}
    />
  );
}
export { Skeleton };

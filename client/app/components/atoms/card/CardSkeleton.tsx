import type React from "react";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { cn } from "~/utils/cn";

interface CardProps {
  className?: string;
}

export function CardSkeleton({
  className,
}: React.PropsWithChildren<CardProps>) {
  return (
    <Skeleton
      className={cn("h-52 rounded-md border-1 border-mauve-6", className)}
    />
  );
}

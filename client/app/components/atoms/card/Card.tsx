import type React from "react";
import { cn } from "~/utils/cn";

interface CardProps {
  className?: string;
}

export function Card({
  className,
  children,
}: React.PropsWithChildren<CardProps>) {
  return (
    <div
      className={cn(
        "h-52 cursor-pointer rounded-md border-1 border-mauve-6 shadow-xs hover:bg-gray-2",
        className,
      )}
    >
      {children}
    </div>
  );
}

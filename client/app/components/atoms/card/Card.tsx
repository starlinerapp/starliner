import React from "react";
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
        "border-mauve-6 hover:bg-gray-2 h-52 cursor-pointer rounded-md border-1 shadow-xs",
        className,
      )}
    >
      {children}
    </div>
  );
}

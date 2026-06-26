import { cn } from "~/utils/cn";

interface LiveIndicatorProps {
  type: "warning" | "success" | "error";
}

export default function LiveIndicator({ type }: LiveIndicatorProps) {
  return (
    <span className="relative flex size-2">
      <span
        className={cn(
          "absolute inline-flex h-full w-full rounded-full opacity-75",
          type === "warning" && "bg-amber-10",
          type === "success" && "bg-grass-10",
          type === "error" && "bg-red-10",
        )}
      ></span>
    </span>
  );
}

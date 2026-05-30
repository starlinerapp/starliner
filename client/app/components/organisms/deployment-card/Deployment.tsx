import React from "react";
import { Play } from "lucide-react";

interface DeploymentTabProps {
  isActive: boolean;
  onSelect: () => void;
}

export function DeploymentTab({ isActive, onSelect }: DeploymentTabProps) {
  return (
    <div className="relative">
      <div className="bg-mauve-8 absolute top-1/2 -left-1 h-2 w-2 -translate-y-1/2 rounded-full" />
      <button
        type="button"
        onClick={onSelect}
        className={
          isActive
            ? "border-violet-9 bg-violet-3 hover:bg-mauve-2 text-violet-9 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border px-4 py-0.5"
            : "border-mauve-6 hover:bg-mauve-2 text-mauve-9 relative z-10 flex cursor-pointer items-center gap-1.5 rounded-md border bg-white px-4 py-0.5"
        }
      >
        <div
          className={
            isActive
              ? "border-violet-9 flex rounded-full border-[1.5px] p-0.5"
              : "border-mauve-9 flex rounded-full border-[1.5px] p-0.5"
          }
        >
          <Play
            className={
              isActive
                ? "fill-violet-9 stroke-violet-9 h-2 w-2"
                : "fill-mauve-9 stroke-mauve-9 h-2 w-2"
            }
          />
        </div>
        Deploy
      </button>
    </div>
  );
}

export function DeploymentLogs() {
  return <div className="text-mauve-11 max-h-125 min-h-24" />;
}

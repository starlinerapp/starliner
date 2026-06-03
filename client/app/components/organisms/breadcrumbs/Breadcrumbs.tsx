import type React from "react";
import { ChevronRight } from "~/components/atoms/icons";

type Crumb = {
  label: React.ReactNode;
};

interface BreadcrumbsProps {
  crumbs: Crumb[];
}

export default function Breadcrumbs({ crumbs }: BreadcrumbsProps) {
  return (
    <div className="border-mauve-6 bg-violet-1 flex h-10 w-full items-center justify-between border-b px-4">
      <span className="flex items-center gap-2 text-sm">
        {crumbs.map((crumb, index) => {
          const isLast = index === crumbs.length - 1;
          return (
            <div key={index} className="flex items-center gap-1">
              <h1 className="text-mauve-12 cursor-default">{crumb.label}</h1>
              {!isLast && <ChevronRight className="h-3 w-3 stroke-2" />}
            </div>
          );
        })}
      </span>
    </div>
  );
}

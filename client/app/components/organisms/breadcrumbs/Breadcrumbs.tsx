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
    <div className="flex h-10 w-full items-center justify-between border-mauve-6 border-b bg-violet-1 px-4">
      <span className="flex items-center gap-2 text-sm">
        {crumbs.map((crumb, index) => {
          const isLast = index === crumbs.length - 1;
          return (
            <div key={index} className="flex items-center gap-1">
              <h1 className="cursor-default text-mauve-12">{crumb.label}</h1>
              {!isLast && <ChevronRight className="h-3 w-3 stroke-2" />}
            </div>
          );
        })}
      </span>
    </div>
  );
}

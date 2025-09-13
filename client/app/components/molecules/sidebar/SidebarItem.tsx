import React from "react";
import { NavLink } from "react-router";
import { cn } from "~/utils/cn";

interface SidebarItemProps {
  title: string;
  icon: React.ReactNode;
  href: string;
}

export default function SidebarItem({ title, icon, href }: SidebarItemProps) {
  return (
    <NavLink className="flex flex-col items-center gap-0.5" to={href}>
      {({ isActive }) => (
        <>
          <div
            className={cn(
              isActive
                ? "text-violet-11 bg-violet-3 border-mauve-6 border-1"
                : "text-violet-12 hover:bg-gray-4 hover:border-gray-4 border-1 border-white",
              "flex h-10 w-10 items-center justify-center rounded-md stroke-1 p-1",
            )}
          >
            {icon}
          </div>
          <div
            className={cn(
              isActive ? "text-violet-11" : "text-violet-12",
              "text-center text-xs",
            )}
          >
            {title}
          </div>
        </>
      )}
    </NavLink>
  );
}

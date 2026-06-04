import type React from "react";
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
                ? "border-1 border-mauve-6 bg-violet-3 text-violet-11"
                : "border-1 border-white text-violet-12 hover:border-gray-4 hover:bg-gray-3",
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

import React from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";
import { NavLink } from "react-router";
import { cn } from "~/utils/cn";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

type SidebarItem = {
  id: string;
  title: string;
  href: string;
};

interface ExtendedSidebarProps {
  title: string;
  sections: SidebarItem[][];
  isLoading: boolean;
  children: React.ReactNode;
}

export default function ExtendedSidebar({
  title,
  sections,
  isLoading,
  children,
}: ExtendedSidebarProps) {
  return (
    <ResizablePanelGroup direction="horizontal">
      <ResizablePanel
        defaultSize={15}
        minSize={10}
        maxSize={20}
        className="bg-violet-1 border-mauve-6 h-screen border-r-1 py-3"
      >
        <div className="text-violet-12 px-4 pb-2 text-sm font-bold">
          {title}
        </div>
        <hr className="border-gray-4 border-t" />
        {isLoading ? (
          <>
            <div className="flex flex-col gap-0.5 p-2">
              <Skeleton className="h-7 w-32" />
            </div>
            <div className="flex flex-col gap-0.5 p-2">
              <Skeleton className="h-7 w-48" />
              <Skeleton className="h-7 w-36" />
              <Skeleton className="h-7 w-52" />
            </div>
          </>
        ) : (
          sections.map((section, i) => (
            <div key={i} className="flex flex-col gap-0.5 p-2">
              {section.map((item) => (
                <NavLink to={item.href} key={item.href} className="flex gap-2">
                  {({ isActive }) => (
                    <span
                      className={cn(
                        "hover:bg-gray-3 flex h-full w-full rounded-md",
                        isActive
                          ? "bg-violet-3 text-violet-11 font-bold"
                          : "text-violet-12",
                      )}
                    >
                      {/* Active Element Indicator */}
                      <span
                        className={cn(
                          isActive && "bg-violet-11 rounded-md",
                          "m-2 w-[3px]",
                        )}
                      />
                      <p className="w-full truncate rounded-md py-2 text-sm">
                        {item.title}
                      </p>
                    </span>
                  )}
                </NavLink>
              ))}
            </div>
          ))
        )}
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={85}>{children}</ResizablePanel>
    </ResizablePanelGroup>
  );
}

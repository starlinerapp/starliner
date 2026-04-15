import React, { useState } from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";
import { NavLink, useLocation } from "react-router";
import { cn } from "~/utils/cn";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { ChevronDown } from "~/components/atoms/icons";
import { motion } from "framer-motion";

type SidebarItem = {
  id: string;
  title: string;
  href: string;
};

type SidebarGroup = {
  id: string;
  title: string;
  children: SidebarItem[];
};

type SidebarSection = SidebarItem[] | SidebarGroup;

interface ExtendedSidebarProps {
  title: string;
  sections: SidebarSection[];
  isLoading: boolean;
  children: React.ReactNode;
}

function SidebarLink({ item }: { item: SidebarItem }) {
  return (
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
          <span
            className={cn(isActive && "bg-violet-11 rounded-md", "m-2 w-[3px]")}
          />
          <p className="w-full truncate rounded-md py-2 text-sm">
            {item.title}
          </p>
        </span>
      )}
    </NavLink>
  );
}

function CollapsibleGroup({ group }: { group: SidebarGroup }) {
  const location = useLocation();
  const isChildActive = group.children.some((child) =>
    location.pathname.startsWith(child.href),
  );
  const [isOpen, setIsOpen] = useState(isChildActive);

  return (
    <div className="flex flex-col gap-0.5">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={cn(
          "hover:bg-gray-3 flex w-full items-center justify-between rounded-md px-3 py-2 text-sm",
          isChildActive ? "text-violet-11 font-bold" : "text-violet-12",
        )}
      >
        <span className="truncate">{group.title}</span>
        <motion.div
          animate={{ rotate: isOpen ? 0 : -90 }}
          transition={{ duration: 0.1, ease: "easeInOut" }}
          className="shrink-0"
        >
          <ChevronDown className="h-4 w-4 stroke-2" />
        </motion.div>
      </button>
      {isOpen && (
        <div className="flex flex-col gap-0.5 pl-3">
          {group.children.map((item) => (
            <SidebarLink key={item.href} item={item} />
          ))}
        </div>
      )}
    </div>
  );
}

function isGroup(section: SidebarSection): section is SidebarGroup {
  return !Array.isArray(section) && "children" in section;
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
            <div className="p-4">
              <Skeleton className="h-5 w-36" />
            </div>
            <div className="flex flex-col gap-3 p-4">
              <Skeleton className="h-5 w-48" />
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-5 w-32" />
            </div>
          </>
        ) : (
          sections.map((section, i) =>
            isGroup(section) ? (
              <div key={section.id} className="p-2">
                <CollapsibleGroup group={section} />
              </div>
            ) : (
              <div key={i} className="flex flex-col gap-0.5 p-2">
                {section.map((item) => (
                  <SidebarLink key={item.href} item={item} />
                ))}
              </div>
            ),
          )
        )}
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={85}>{children}</ResizablePanel>
    </ResizablePanelGroup>
  );
}

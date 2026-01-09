import React from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";
import ArchitectureCanvas from "~/components/organisms/canvas/ArchitectureCanvas";
import { Outlet, useParams } from "react-router";
import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";

export default function Layout() {
  const { slug, id, environment } = useParams<{
    slug: string;
    id: string;
    environment: string;
  }>();

  const navigationBarItems = [
    {
      title: "Git Repository",
      href: `/${slug}/projects/${id}/${environment}/architecture/git`,
    },
    {
      title: "Database",
      href: `/${slug}/projects/${id}/${environment}/architecture/database`,
    },
  ];

  return (
    <ResizablePanelGroup direction="horizontal" className="h-full">
      <ResizablePanel
        defaultSize={70}
        className="border-mauve-6 h-full border-r-1"
      >
        <ArchitectureCanvas />
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={30} className="flex h-full flex-col">
        <LinkNavigationBar items={navigationBarItems} />
        <div className="p-4">
          <Outlet />
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

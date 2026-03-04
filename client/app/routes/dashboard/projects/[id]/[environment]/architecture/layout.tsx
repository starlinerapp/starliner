import React, { useMemo } from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";
import ArchitectureCanvas from "~/components/organisms/canvas/ArchitectureCanvas";
import { Outlet, useOutletContext, useParams } from "react-router";
import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import type { ResponseEnvironment } from "~/server/api/client/generated";
import { ReactFlowProvider } from "@xyflow/react";

type ContextType = {
  environment: ResponseEnvironment;
  clusterId: number | undefined;
};

export default function Layout() {
  const trpc = useTRPC();
  const { slug, id, environment } = useParams<{
    slug: string;
    id: string;
    environment: string;
  }>();

  const { data: project } = useQuery(
    trpc.project.getProject.queryOptions({ id: Number(id) }),
  );
  const currentEnvironment = useMemo(
    () => project?.environments.find((e) => e.slug === environment),
    [project, environment],
  );

  const navigationBarItems = [
    {
      title: "Git Repository",
      href: `/${slug}/projects/${id}/${environment}/architecture/git`,
    },
    {
      title: "Image",
      href: `/${slug}/projects/${id}/${environment}/architecture/image`,
    },
    {
      title: "Ingress",
      href: `/${slug}/projects/${id}/${environment}/architecture/ingress`,
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
        {currentEnvironment && (
          <ReactFlowProvider>
            <ArchitectureCanvas environment={currentEnvironment} />
          </ReactFlowProvider>
        )}
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={30} className="flex h-full flex-col">
        <div className="bg-violet-1">
          <LinkNavigationBar items={navigationBarItems} />
        </div>
        <div className="p-4">
          <Outlet
            context={{
              environment: currentEnvironment,
              clusterId: project?.clusterId,
            }}
          />
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

export function useEnvironment() {
  return useOutletContext<ContextType>();
}

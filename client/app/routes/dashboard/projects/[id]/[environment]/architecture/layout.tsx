import React, { useEffect, useMemo, useRef } from "react";
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
import BottomBar from "~/components/organisms/bottom-bar/deployment/BottomBar";
import type { ImperativePanelHandle } from "react-resizable-panels";

type ContextType = {
  environment: ResponseEnvironment;
  clusterId: number | undefined;
  teamId: number | undefined;
};

export default function Layout() {
  const trpc = useTRPC();
  const { slug, id, environment, deploymentId } = useParams<{
    slug: string;
    id: string;
    environment: string;
    deploymentId: string;
  }>();

  const { data: project } = useQuery(
    trpc.project.getProject.queryOptions({ id: Number(id) }),
  );

  const currentEnvironment = useMemo(
    () => project?.environments.find((e) => e.slug === environment),
    [project, environment],
  );

  const { data: environmentDeployments } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: Number(currentEnvironment?.id) },
      { enabled: !!currentEnvironment },
    ),
  );

  const currentDeployment = useMemo(() => {
    if (!environmentDeployments) return undefined;
    const allDeployments = [
      ...environmentDeployments.gitDeployments,
      ...environmentDeployments.images,
      ...environmentDeployments.ingresses,
      ...environmentDeployments.databases,
    ];

    return allDeployments.find((d) => d.id === Number(deploymentId));
  }, [environmentDeployments, deploymentId]);

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

  const bottomPanelRef = useRef<ImperativePanelHandle>(null);
  useEffect(() => {
    if (!deploymentId) return;

    const panel = bottomPanelRef.current;
    if (!panel) return;

    const currentSize = panel.getSize();

    if (currentSize < 10) {
      panel.resize(45);
    }
  }, [deploymentId]);

  return (
    <ResizablePanelGroup direction="horizontal" className="h-full w-full">
      <ResizablePanel defaultSize={70} className="h-full">
        <ResizablePanelGroup
          direction="vertical"
          className="border-mauve-6 h-full w-full border-r-1"
        >
          <ResizablePanel defaultSize={70} className="h-full">
            {currentEnvironment && (
              <ReactFlowProvider>
                <ArchitectureCanvas environment={currentEnvironment} />
              </ReactFlowProvider>
            )}
          </ResizablePanel>

          <ResizableHandle />

          <ResizablePanel
            minSize={4}
            maxSize={85}
            ref={bottomPanelRef}
            defaultSize={3}
            className="border-mauve-6 border-t-1"
          >
            <BottomBar deployment={currentDeployment} />
          </ResizablePanel>
        </ResizablePanelGroup>
      </ResizablePanel>

      <ResizableHandle />

      <ResizablePanel defaultSize={30} className="flex h-full flex-col">
        <div className="bg-violet-1">
          <LinkNavigationBar items={navigationBarItems} />
        </div>
        <div className="max-h-[calc(100vh-135px)] flex-1 overflow-auto p-4">
          <Outlet
            context={{
              environment: currentEnvironment,
              clusterId: project?.clusterId,
              teamId: project?.teamId,
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

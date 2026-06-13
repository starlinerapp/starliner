import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { useNavigate, useParams } from "react-router";
import Button from "~/components/atoms/button/Button";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { Download } from "~/components/atoms/icons";
import LiveIndicator from "~/components/atoms/live-indicator/LiveIndicator";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import BottomBar from "~/components/organisms/bottom-bar/cluster/BottomBar";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

export default function General() {
  const navigate = useNavigate();
  const { slug, id } = useParams<{ slug: string; id: string }>();

  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const {
    data: clusterData,
    isLoading,
    error,
  } = useQuery(
    trpc.cluster.getCluster.queryOptions(
      { id: Number(id) },
      { refetchInterval: 4000 },
    ),
  );

  useEffect(() => {
    if (error) {
      (async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationClusters.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${slug}/clusters/all`);
      })();
    }
  }, [
    error,
    slug,
    organization.id,
    trpc.organization.getOrganizationClusters.queryKey,
    queryClient.invalidateQueries,
    navigate,
  ]);

  type LiveIndicatorType = "warning" | "success" | "error";
  const liveIndicatorStatusMap: Record<string, LiveIndicatorType> = {
    pending: "warning",
    running: "success",
    deleted: "error",
  };
  const liveIndicatorType =
    liveIndicatorStatusMap[clusterData?.status ?? "pending"];

  const statusMap = {
    pending: "Creating",
    running: "Live",
    deleted: "Deleting",
  };
  const status = statusMap[clusterData?.status ?? "pending"];
  const showBottomPanel = Boolean(
    clusterData?.id &&
      (clusterData.status === "pending" || clusterData.status === "running"),
  );

  return (
    <ResizablePanelGroup
      direction="vertical"
      className="h-full w-full [&:not(:has([data-resize-handle-state=drag]))_[data-slot=resizable-panel]]:transition-[flex-grow,flex-basis] [&:not(:has([data-resize-handle-state=drag]))_[data-slot=resizable-panel]]:duration-200 [&:not(:has([data-resize-handle-state=drag]))_[data-slot=resizable-panel]]:ease-in-out"
    >
      <ResizablePanel
        defaultSize={showBottomPanel ? 55 : 100}
        className="h-full overflow-auto"
      >
        <div className="w-full p-4">
          <div className="rounded-md border border-mauve-6 bg-gray-2 text-sm shadow-xs">
            <div className="flex h-14 items-center rounded-t-md px-4 font-bold text-mauve-12 text-xs uppercase">
              Details
            </div>
            <div className="mx-1 mb-1 divide-y divide-mauve-6 overflow-hidden rounded-md border border-mauve-6 bg-white shadow-xs">
              <div className="flex items-center justify-between gap-2 px-4 py-2">
                <div>
                  <h2 className="text-mauve-12">Status</h2>
                </div>
                {isLoading ? (
                  <Skeleton className="h-5 w-24" />
                ) : (
                  <span className="flex items-center gap-3">
                    <LiveIndicator type={liveIndicatorType} />
                    <p className="pr-2 text-mauve-11 capitalize">{status}</p>
                  </span>
                )}
              </div>
              <div className="flex items-center justify-between gap-2 px-4 py-2">
                <div>
                  <h2 className="text-mauve-12">Server Type</h2>
                </div>
                {isLoading ? (
                  <Skeleton className="h-5 w-24" />
                ) : (
                  <CopyToClipboard
                    className="text-mauve-11"
                    text={clusterData?.serverType ?? ""}
                  />
                )}
              </div>
              <div className="flex items-center justify-between gap-2 px-4 py-2">
                <div>
                  <h2 className="text-mauve-12">IPv4 Address</h2>
                </div>
                {isLoading || clusterData?.status === "pending" ? (
                  <Skeleton className="h-5 w-32" />
                ) : (
                  <CopyToClipboard
                    className="text-mauve-11"
                    text={clusterData?.ipv4Address ?? ""}
                  />
                )}
              </div>
              <div className="flex items-center justify-between gap-2 px-4 py-2">
                <div>
                  <h2 className="text-mauve-12">User</h2>
                </div>
                {isLoading || clusterData?.status === "pending" ? (
                  <Skeleton className="h-5 w-14" />
                ) : (
                  <CopyToClipboard
                    className="text-mauve-11"
                    text={clusterData?.user ?? ""}
                  />
                )}
              </div>
              <div className="flex items-center justify-between gap-2 px-4 py-2">
                <div>
                  <h2 className="text-mauve-12">SSH Key</h2>
                  <p className="text-mauve-11 text-xs">
                    You can use the SSH key to access the cluster.
                  </p>
                </div>
                {isLoading || clusterData?.status === "pending" ? (
                  <Skeleton className="h-9.5 w-32" />
                ) : (
                  <a
                    href={`/clusters/${id}/resources/private-key`}
                    download="private-key.pem"
                  >
                    <Button intent="secondary" className="w-32">
                      <Download width={18} strokeWidth={2} />
                      Download
                    </Button>
                  </a>
                )}
              </div>
            </div>
          </div>
        </div>
      </ResizablePanel>

      {showBottomPanel && clusterData && (
        <>
          <ResizableHandle />

          <ResizablePanel
            defaultSize={45}
            minSize={4}
            maxSize={85}
            className="border-mauve-6 border-t"
          >
            <BottomBar clusterId={clusterData.id} status={clusterData.status} />
          </ResizablePanel>
        </>
      )}
    </ResizablePanelGroup>
  );
}

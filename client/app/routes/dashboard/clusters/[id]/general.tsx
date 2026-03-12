import React, { useEffect } from "react";
import Button from "~/components/atoms/button/Button";
import { Download } from "~/components/atoms/icons";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { useNavigate, useParams } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import LiveIndicator from "~/components/atoms/live-indicator/LiveIndicator";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

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
      {
        refetchInterval: 4000,
      },
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
  }, [error, id, slug]);

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

  return (
    <div className="w-full p-4 xl:w-3/5">
      <div className="border-mauve-6 rounded-md border-1 text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
          Details
        </div>
        <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
          <div>
            <h1 className="text-mauve-12">Status</h1>
          </div>
          {isLoading ? (
            <Skeleton className="h-5 w-24" />
          ) : (
            <span className="flex items-center gap-3">
              <LiveIndicator type={liveIndicatorType} />
              <p className="text-mauve-11 pr-2 capitalize">{status}</p>
            </span>
          )}
        </div>
        <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
          <div>
            <h1 className="text-mauve-12">Server Type</h1>
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
        <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
          <div>
            <h1 className="text-mauve-12">IPv4 Address</h1>
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
        <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
          <div>
            <h1 className="text-mauve-12">User</h1>
          </div>
          {isLoading || clusterData?.status === "pending" ? (
            <Skeleton className="h-5 w-14" />
          ) : (
            <CopyToClipboard className="text-mauve-11" text="root" />
          )}
        </div>
        <div className="flex items-center justify-between px-4 py-2">
          <div>
            <h1 className="text-mauve-12">SSH Key</h1>
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
  );
}

import React from "react";
import Button from "~/components/atoms/button/Button";
import { Download } from "~/components/atoms/icons";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { useParams } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";

export default function General() {
  const { id } = useParams<{ id: string }>();

  const trpc = useTRPC();
  const { data: clusterData } = useQuery(
    trpc.cluster.getCluster.queryOptions({ id: Number(id) }),
  );

  return (
    <div className="w-3/5 p-4">
      <div className="border-mauve-6 rounded-md border-1">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-sm uppercase">
          Details
        </div>
        <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
          <div>
            <h1 className="text-mauve-12">IPv4 Address</h1>
          </div>
          <CopyToClipboard
            className="text-mauve-11"
            text={clusterData?.ipv4Address ?? ""}
          />
        </div>
        <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
          <div>
            <h1 className="text-mauve-12">User</h1>
          </div>
          <CopyToClipboard className="text-mauve-11" text="root" />
        </div>
        <div className="flex items-center justify-between px-4 py-2">
          <div>
            <h1 className="text-mauve-12">SSH Key</h1>
            <p className="text-mauve-11 text-sm">
              You can use the SSH key to access the cluster.
            </p>
          </div>
          <a
            href={`/clusters/${id}/resources/private-key`}
            download="private-key.pem"
          >
            <Button intent="secondary" className="w-32">
              <Download width={18} strokeWidth={2} />
              Download
            </Button>
          </a>
        </div>
      </div>
    </div>
  );
}

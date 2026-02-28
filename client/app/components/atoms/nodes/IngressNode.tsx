import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { EllipsisVertical, Shuffle, Trash } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import React from "react";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import * as Popover from "@radix-ui/react-popover";

type IngressNode = Node<{
  id: number;
  serviceName: string;
  status: string;
  port: string;
}>;

export default function IngressNode({ data }: NodeProps<IngressNode>) {
  return (
    <div className="bg-white-a12 text-mauve-11">
      <Handle
        type="target"
        position={Position.Right}
        className="!border-mauve-8 !h-3 !w-3 !border-1 !bg-white"
      />
      <div className="database-node border-mauve-6 bg-mauve-2 flex w-[350px] flex-col gap-2 rounded-md border-1 p-2 shadow-md">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <Shuffle className="w-5" />
            <p>{data.serviceName}</p>
          </div>
          <IngressContextMenu deploymentId={data.id} />
        </div>
        <div>
          <div className="bg-gray-2 border-mauve-6 flex justify-between rounded-t-md border-1 p-2 text-sm shadow-md">
            <p>Status</p>
            <span className="flex items-center gap-1.5">
              <span
                className={cn(
                  "h-3 w-3 rounded-full",
                  data.status === "healthy" ? "bg-grass-9" : "bg-red-9",
                )}
              ></span>
              <p>{data.status}</p>
            </span>
          </div>
          <div className="bg-white-a12 border-mauve-6 -mt-1.5 flex flex-col gap-2 rounded-md border-1 p-2 text-sm shadow-sm">
            <span className="flex justify-between">
              <p>Port</p>
              <CopyToClipboard className="text-mauve-11" text={data.port} />
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

interface IngressContextMenuProps {
  deploymentId: number;
}

function IngressContextMenu({ deploymentId }: IngressContextMenuProps) {
  const trpc = useTRPC();
  const deleteDeploymentMutation = useMutation(
    trpc.deployment.deleteDeployment.mutationOptions(),
  );

  function handleDeleteClicked() {
    deleteDeploymentMutation.mutate({
      id: deploymentId,
    });
  }

  return (
    <Popover.Root>
      <Popover.Trigger className="hover:bg-gray-4 flex h-7 w-7 cursor-pointer rounded-md p-1">
        <EllipsisVertical className="w-6" />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          side="bottom"
          align="start"
          className="border-gray-6 m-2 rounded-md border bg-white shadow-md"
        >
          <div className="flex min-w-[120px] flex-col p-0.5">
            <Popover.Close asChild>
              <button
                className="hover:bg-gray-3 text-mauve-11 flex w-full cursor-pointer flex-row items-center gap-2 rounded-md p-2 text-sm"
                onClick={handleDeleteClicked}
              >
                <Trash className="w-5" />
                <p>Delete</p>
              </button>
            </Popover.Close>
          </div>
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  );
}

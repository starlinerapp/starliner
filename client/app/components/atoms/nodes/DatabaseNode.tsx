import React from "react";
import * as Popover from "@radix-ui/react-popover";
import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { Database, EllipsisVertical, Trash } from "~/components/atoms/icons";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { cn } from "~/utils/cn";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";

type DatabaseNode = Node<{
  id: number;
  serviceName: string;
  status: string;
  port: string;
  username: string;
  password: string;
}>;

export default function DatabaseNode({ data }: NodeProps<DatabaseNode>) {
  return (
    <div className="bg-white-a12 text-mauve-11">
      <Handle
        type="target"
        position={Position.Left}
        className="!border-mauve-8 !h-3 !w-3 !border-1 !bg-white"
      />
      <div className="database-node border-mauve-6 bg-mauve-2 flex w-[350px] flex-col gap-2 rounded-md border-1 p-2 shadow-md">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <Database className="w-5" />
            <p>{data.serviceName}</p>
          </div>
          <DatabaseContextMenu deploymentId={data.id} />
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
              <p>Username</p>
              <CopyToClipboard className="text-mauve-11" text={data.username} />
            </span>
            <span className="flex justify-between">
              <p>Password</p>
              <CopyToClipboard
                masked={true}
                className="text-mauve-11"
                text={data.password}
              />
            </span>
            <span className="flex justify-between">
              <p>Internal Endpoint</p>
              <CopyToClipboard
                className="text-mauve-11"
                text={`${data.serviceName}-rw:${data.port}`}
              />
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

interface DatabaseContextMenuProps {
  deploymentId: number;
}

function DatabaseContextMenu({ deploymentId }: DatabaseContextMenuProps) {
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
      <Popover.Trigger
        className="hover:bg-gray-4 flex h-7 w-7 cursor-pointer rounded-md p-1"
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
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
                onClick={(e) => {
                  e.stopPropagation();
                  handleDeleteClicked();
                }}
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

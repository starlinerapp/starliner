import React from "react";
import * as Popover from "@radix-ui/react-popover";
import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { cn } from "~/utils/cn";
import { EllipsisVertical, GitBranch, Trash } from "~/components/atoms/icons";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useLocation, useMatch, useNavigate } from "react-router";

type GitNode = Node<{
  id: number;
  serviceName: string;
  internalEndpoint: string;
  status: string;
  port: string;
  gitUrl: string;
}>;

export default function GitNode({ data, selected }: NodeProps<GitNode>) {
  return (
    <div
      className={cn(
        "bg-white-a12 text-mauve-11 hover:ring-violet-6 hover:rounded-md hover:ring-2",
        selected && "ring-violet-8 hover:ring-violet-8 rounded-md ring-2",
      )}
    >
      <Handle
        type="target"
        isConnectable={false}
        position={Position.Left}
        className="!border-mauve-8 !h-3 !w-3 !border-1 !bg-white"
      />
      <Handle
        type="source"
        isConnectable={false}
        position={Position.Right}
        className="!border-mauve-8 !h-3 !w-3 !border-1 !bg-white"
      />
      <div className="database-node border-mauve-6 bg-mauve-2 flex w-[350px] flex-col gap-2 rounded-md border-1 p-2 shadow-md">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <GitBranch className="w-5" />
            <p>{data.serviceName}</p>
          </div>
          <GitContextMenu deploymentId={data.id} />
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
            <span className="flex min-w-0 items-center gap-2">
              <p className="shrink-0">Repository</p>
              <span className="min-w-0 flex-1 overflow-hidden [&>*]:block [&>*]:max-w-full [&>*]:truncate">
                <CopyToClipboard
                  className="text-mauve-11 block max-w-full truncate"
                  text={data.gitUrl}
                />
              </span>
            </span>
            <span className="flex justify-between">
              <p>Branch</p>
              <CopyToClipboard className="text-mauve-11" text="main" />
            </span>
            <span className="flex justify-between">
              <p>Internal Endpoint</p>
              <CopyToClipboard
                className="text-mauve-11"
                text={data.internalEndpoint}
              />
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

interface GitContextMenuProps {
  deploymentId: number;
}

function GitContextMenu({ deploymentId }: GitContextMenuProps) {
  const trpc = useTRPC();
  const deleteDeploymentMutation = useMutation(
    trpc.deployment.deleteDeployment.mutationOptions(),
  );

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const isOnDatabaseDetail = !!useMatch(
    "/starliner/projects/:id/:environment/architecture/git/:deploymentId",
  );

  function handleDeleteClicked() {
    deleteDeploymentMutation.mutate(
      {
        id: deploymentId,
      },
      {
        onSuccess: () => {
          if (isOnDatabaseDetail) {
            const parent = pathname.replace(/\/[^/]+\/?$/, "");
            navigate(parent, { relative: "path", replace: true });
          }
        },
      },
    );
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

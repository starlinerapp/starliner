import { useMutation } from "@tanstack/react-query";
import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { useState } from "react";
import { useLocation, useMatch, useNavigate } from "react-router";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { EllipsisVertical, GitBranch, Trash } from "~/components/atoms/icons";
import {
  Popover,
  PopoverClose,
  PopoverContent,
  PopoverTrigger,
} from "~/components/atoms/popover/Popover";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

type GitNode = Node<{
  id: number;
  serviceName: string;
  internalEndpoint: string;
  status: string;
  port: string;
  gitUrl: string;
}>;

export default function GitNode({ data, selected }: NodeProps<GitNode>) {
  const [isDeleting, setIsDeleting] = useState(false);

  return (
    <div
      className={cn(
        "bg-white-a12 text-mauve-11 hover:rounded-md hover:ring-2 hover:ring-violet-6",
        selected && "rounded-md ring-2 ring-violet-8 hover:ring-violet-8",
        isDeleting && "pointer-events-none grayscale",
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
      <div className="database-node flex w-[350px] flex-col gap-2 rounded-md border-1 border-mauve-6 bg-mauve-2 p-2 shadow-md">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <GitBranch className="w-5" />
            <p>{data.serviceName}</p>
          </div>
          <GitContextMenu
            deploymentId={data.id}
            setIsDeleting={setIsDeleting}
          />
        </div>
        <div>
          <div className="flex justify-between rounded-t-md border-1 border-mauve-6 bg-gray-2 p-2 text-sm shadow-md">
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
          <div className="-mt-1.5 flex flex-col gap-2 rounded-md border-1 border-mauve-6 bg-white-a12 p-2 text-sm shadow-sm">
            <span className="flex min-w-0 items-center gap-2">
              <p className="shrink-0">Repository</p>
              <span className="min-w-0 flex-1 overflow-hidden [&>*]:block [&>*]:max-w-full [&>*]:truncate">
                <CopyToClipboard
                  className="block max-w-full truncate text-mauve-11"
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
  setIsDeleting: (v: boolean) => void;
}

function GitContextMenu({ deploymentId, setIsDeleting }: GitContextMenuProps) {
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
    setIsDeleting(true);
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
        onError: () => {
          setIsDeleting(false);
        },
      },
    );
  }

  return (
    <Popover>
      <PopoverTrigger
        className="flex h-7 w-7 cursor-pointer rounded-md p-1 hover:bg-gray-4"
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        <EllipsisVertical className="w-6" />
      </PopoverTrigger>
      <PopoverContent side="bottom" align="start" sideOffset={4}>
        <div className="flex min-w-[120px] flex-col p-0.5">
          <PopoverClose asChild>
            <button
              type="button"
              className="flex w-full cursor-pointer flex-row items-center gap-2 rounded-md p-2 text-mauve-11 text-sm hover:bg-gray-3"
              onClick={(e) => {
                e.stopPropagation();
                handleDeleteClicked();
              }}
            >
              <Trash className="w-5" />
              <p>Delete</p>
            </button>
          </PopoverClose>
        </div>
      </PopoverContent>
    </Popover>
  );
}

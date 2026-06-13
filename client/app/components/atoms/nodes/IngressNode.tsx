import { useMutation } from "@tanstack/react-query";
import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { useState } from "react";
import { useLocation, useMatch, useNavigate } from "react-router";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { EllipsisVertical, Shuffle, Trash } from "~/components/atoms/icons";
import {
  Popover,
  PopoverClose,
  PopoverContent,
  PopoverTrigger,
} from "~/components/atoms/popover/Popover";
import type { ResponseIngressHost } from "~/server/api/clients/server/generated";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

type IngressNode = Node<{
  id: number;
  serviceName: string;
  status: string;
  port: string;
  hosts: ResponseIngressHost[];
}>;

export default function IngressNode({
  data,
  selected,
}: NodeProps<IngressNode>) {
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
        type="source"
        isConnectable={false}
        position={Position.Right}
        className="!border-mauve-8 !h-3 !w-3 !border-1 !bg-white"
      />
      <div className="database-node flex w-[350px] flex-col gap-2 rounded-md border-1 border-mauve-6 bg-mauve-2 p-2 shadow-md">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <Shuffle className="w-5" />
            <p>{data.serviceName}</p>
          </div>
          <IngressContextMenu
            deploymentId={data.id}
            setIsDeleting={setIsDeleting}
          />
        </div>
        <div>
          <div className="-mt-1.5 flex flex-col gap-2 rounded-md border-1 border-mauve-6 bg-white-a12 p-2 text-sm shadow-sm">
            <div className="flex flex-col">
              <p>Ingress</p>
              {data.hosts.map((host, hostIndex) => {
                return (
                  <div key={host.host} className="relative">
                    <div
                      className={cn(
                        "relative flex flex-col gap-1 border-mauve-6 border-l-2 pl-6",
                        hostIndex === 0 ? "pt-1" : "pt-2",
                      )}
                    >
                      <div className="absolute -left-[2px] h-3 w-5 rounded-bl-md border-mauve-6 border-b-2 border-l-2" />
                      <span className="min-w-0 flex-1 overflow-hidden [&>*]:block [&>*]:max-w-full [&>*]:truncate">
                        <CopyToClipboard
                          className="block max-w-full truncate px-1 font-medium text-mauve-11 text-sm"
                          text={host.host}
                        />
                      </span>
                      <div className="relative flex flex-col gap-2 border-mauve-6 border-l-2 pl-6">
                        {host.paths?.map((path) => {
                          return (
                            <div
                              key={path.path}
                              className="relative flex flex-col gap-1"
                            >
                              <div className="absolute -left-6.5 h-3 w-5 rounded-bl-md border-mauve-6 border-b-2 border-l-2" />
                              <div className="flex items-center gap-2 text-sm">
                                <div className="flex min-w-0 items-center gap-2">
                                  {path.pathType ? (
                                    <span className="shrink-0 rounded border border-indigo-6 bg-violet-2 px-1.5 py-0.5 font-medium font-mono text-[10px] text-violet-11 uppercase tracking-wide">
                                      {path.pathType}
                                    </span>
                                  ) : null}
                                  <span className="truncate">{path.path}</span>
                                </div>

                                <div className="flex flex-1 items-center">
                                  <div className="h-[2px] flex-1 bg-mauve-6" />
                                  <div className="h-0 w-0 border-mauve-6 border-y-5 border-y-transparent border-l-12" />
                                </div>

                                <span className="shrink-0 rounded-md bg-mauve-3 px-2 py-0.5 font-medium text-mauve-11 text-xs ring-1 ring-mauve-5 ring-inset">
                                  {path.serviceName}
                                </span>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

interface IngressContextMenuProps {
  deploymentId: number;
  setIsDeleting: (v: boolean) => void;
}

function IngressContextMenu({
  deploymentId,
  setIsDeleting,
}: IngressContextMenuProps) {
  const trpc = useTRPC();
  const deleteDeploymentMutation = useMutation(
    trpc.deployment.deleteDeployment.mutationOptions(),
  );

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const isOnIngressDetail = !!useMatch(
    "/starliner/projects/:id/:environment/architecture/ingress/:deploymentId",
  );

  function handleDeleteClicked() {
    setIsDeleting(true);

    deleteDeploymentMutation.mutate(
      {
        id: deploymentId,
      },
      {
        onSuccess: () => {
          if (isOnIngressDetail) {
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

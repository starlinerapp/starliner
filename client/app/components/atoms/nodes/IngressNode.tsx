import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { EllipsisVertical, Shuffle, Trash } from "~/components/atoms/icons";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import React from "react";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import * as Popover from "@radix-ui/react-popover";
import type { ResponseIngressHost } from "~/server/api/client/generated";
import { cn } from "~/utils/cn";
import { useLocation, useMatch, useNavigate } from "react-router";

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
  return (
    <div
      className={cn(
        "bg-white-a12 text-mauve-11 hover:ring-violet-6 hover:rounded-md hover:ring-2",
        selected && "ring-violet-8 hover:ring-violet-8 rounded-md ring-2",
      )}
    >
      <Handle
        type="source"
        isConnectable={false}
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
          <div className="bg-white-a12 border-mauve-6 -mt-1.5 flex flex-col gap-2 rounded-md border-1 p-2 text-sm shadow-sm">
            <div className="flex flex-col">
              <p>Ingress</p>
              {data.hosts.map((host, hostIndex) => {
                return (
                  <div key={hostIndex} className="relative">
                    <div
                      className={cn(
                        "border-mauve-6 relative flex flex-col gap-1 border-l-2 pl-6",
                        hostIndex === 0 ? "pt-1" : "pt-2",
                      )}
                    >
                      <div className="border-mauve-6 absolute -left-[2px] h-3 w-5 rounded-bl-md border-b-2 border-l-2" />
                      <div>
                        <CopyToClipboard
                          className="text-mauve-11 px-1 text-sm font-medium"
                          text={host.host}
                        />
                      </div>
                      <div className="border-mauve-6 relative flex flex-col gap-2 border-l-2 pl-6">
                        {host.paths?.map((path, pathIndex) => {
                          return (
                            <div
                              key={pathIndex}
                              className="relative flex flex-col gap-1"
                            >
                              <div className="border-mauve-6 absolute -left-6.5 h-3 w-5 rounded-bl-md border-b-2 border-l-2" />
                              <div className="flex items-center gap-2 text-sm">
                                <div className="flex min-w-0 items-center gap-2">
                                  {path.pathType ? (
                                    <span className="border-indigo-6 bg-violet-2 text-violet-11 shrink-0 rounded border px-1.5 py-0.5 font-mono text-[10px] font-medium tracking-wide uppercase">
                                      {path.pathType}
                                    </span>
                                  ) : null}
                                  <span className="truncate">{path.path}</span>
                                </div>

                                <div className="flex flex-1 items-center">
                                  <div className="bg-mauve-6 h-[2px] flex-1" />
                                  <div className="border-mauve-6 h-0 w-0 border-y-5 border-l-12 border-y-transparent" />
                                </div>

                                <span className="bg-mauve-3 text-mauve-11 ring-mauve-5 shrink-0 rounded-md px-2 py-0.5 text-xs font-medium ring-1 ring-inset">
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
}

function IngressContextMenu({ deploymentId }: IngressContextMenuProps) {
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

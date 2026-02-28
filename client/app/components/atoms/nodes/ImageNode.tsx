import { Handle, type Node, type NodeProps, Position } from "@xyflow/react";
import { CodeBracket } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import React from "react";

type ImageNode = Node<{
  id: number;
  serviceName: string;
  status: string;
  port: string;
  imageName: string;
  tag: string;
}>;

export default function ImageNode({ data }: NodeProps<ImageNode>) {
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
            <CodeBracket className="w-5" />
            <p>{data.serviceName}</p>
          </div>
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
              <p>Image</p>
              <CopyToClipboard
                className="text-mauve-11"
                text={data.imageName}
              />
            </span>
            <span className="flex justify-between">
              <p>Tag</p>
              <CopyToClipboard className="text-mauve-11" text={data.tag} />
            </span>
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

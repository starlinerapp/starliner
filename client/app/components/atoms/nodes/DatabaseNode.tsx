import React from "react";
import * as Popover from "@radix-ui/react-popover";
import type { Node, NodeProps } from "@xyflow/react";
import { Database, EllipsisVertical, Trash } from "~/components/atoms/icons";

type DatabaseNode = Node<{
  id: string;
  serviceName: string;
  port: string;
  username: string;
  password: string;
}>;

export default function DatabaseNode({ data }: NodeProps<DatabaseNode>) {
  return (
    <div className="bg-white-a12 text-mauve-11">
      <div className="database-node border-mauve-6 bg-mauve-2 flex w-[300px] flex-col gap-2 rounded-md border-1 p-2">
        <div className="flex justify-between">
          <div className="flex items-center gap-2">
            <Database className="w-5" />
            <p>{data.serviceName}</p>
          </div>
          <DatabaseContextMenu deploymentId={data.id} />
        </div>
        <div className="bg-white-a12 border-mauve-6 flex flex-col gap-2 rounded-md border-1 p-2 text-sm">
          <span className="flex justify-between">
            <p>Port</p>
            <p>{data.port}</p>
          </span>
          <span className="flex justify-between">
            <p>Username</p>
            <p>{data.username}</p>
          </span>
          <span className="flex justify-between">
            <p>Password</p>
            <p>{data.password}</p>
          </span>
        </div>
      </div>
    </div>
  );
}

interface DatabaseContextMenuProps {
  deploymentId: string;
}

function DatabaseContextMenu({ deploymentId }: DatabaseContextMenuProps) {
  function handleDeleteClicked() {
    console.log("Delete clicked - ", deploymentId);
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

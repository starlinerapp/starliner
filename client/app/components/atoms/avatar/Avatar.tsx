import React from "react";
import { useQuery } from "@tanstack/react-query";
import * as Popover from "@radix-ui/react-popover";
import { useNavigate } from "react-router";
import { authClient } from "~/utils/auth/client";

interface AvatarIconProps {
  name: string;
}

function AvatarIcon({ name }: AvatarIconProps) {
  return (
    <div className="bg-violet-9 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-sm text-white">
      {name.substring(0, 2).toUpperCase()}
    </div>
  );
}

export default function Avatar() {
  const navigate = useNavigate();
  const { data: session } = useQuery({
    queryFn: () => authClient.getSession(),
    queryKey: ["session"],
  });

  const username = session?.data?.user.name ?? "";

  async function handleSignOutClicked() {
    await authClient.signOut();
    navigate("/login");
  }

  return (
    <Popover.Root>
      <Popover.Trigger className="data-[state=open]:bg-gray-4 data-[state=open]:border-gray-4 hover:bg-gray-4 hover:border-gray-4 flex h-full w-full items-center justify-center rounded-md border-1 border-white">
        <AvatarIcon name={username} />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          side="right"
          align="start"
          className="border-gray-6 m-2 rounded-md border bg-white shadow-md"
        >
          <div className="flex min-w-[175px] flex-col p-1">
            <div className="flex gap-2 p-1">
              <AvatarIcon name={username} />
              <div className="flex flex-col">
                <p className="text-gray-12 text-xs font-bold">
                  {session?.data?.user.name}
                </p>
                <p className="text-gray-11 text-xs">
                  {session?.data?.user.email}
                </p>
              </div>
            </div>
            <button
              className="hover:bg-gray-3 flex flex-row items-center gap-2 rounded-md p-2 text-xs"
              onClick={async () => {}}
            >
              <p>User Settings</p>
            </button>
            <button
              className="hover:bg-gray-3 flex flex-row items-center gap-2 rounded-md p-2 text-xs"
              onClick={handleSignOutClicked}
            >
              <p>Sign Out</p>
            </button>
          </div>
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  );
}

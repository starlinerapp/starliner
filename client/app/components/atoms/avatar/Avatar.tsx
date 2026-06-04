import * as Popover from "@radix-ui/react-popover";
import { useNavigate } from "react-router";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { getAuthClient } from "~/utils/auth/client";

interface AvatarIconProps {
  name: string;
  profilePicture: string | undefined | null;
}

export function AvatarIcon({ name, profilePicture }: AvatarIconProps) {
  return (
    <>
      {profilePicture ? (
        <div className="h-8 w-8 overflow-hidden rounded-full">
          <img
            src={profilePicture}
            alt="User avatar"
            className="h-full w-full object-cover"
          />
        </div>
      ) : (
        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-violet-9 text-sm text-white">
          {name.substring(0, 2).toUpperCase()}
        </div>
      )}
    </>
  );
}

export default function Avatar() {
  const authClient = getAuthClient();
  const navigate = useNavigate();
  const { data: session, isPending: isSessionPending } =
    authClient.useSession();

  const username = session?.user.name ?? "";
  const profilePicture = session?.user?.image;

  async function handleSignOutClicked() {
    await authClient.signOut();
    navigate("/login");
  }

  return (
    <Popover.Root>
      <Popover.Trigger className="flex h-full w-full cursor-pointer items-center justify-center rounded-md border-1 border-white hover:border-gray-4 hover:bg-gray-4 data-[state=open]:border-gray-4 data-[state=open]:bg-gray-4">
        {isSessionPending ? (
          <Skeleton className="h-8 w-8 rounded-full" />
        ) : (
          <AvatarIcon name={username} profilePicture={profilePicture} />
        )}
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          side="right"
          align="start"
          className="m-2 rounded-md border border-gray-6 bg-white shadow-md"
        >
          <div className="flex min-w-[175px] flex-col p-1">
            <div className="flex gap-2 p-1">
              <AvatarIcon name={username} profilePicture={profilePicture} />
              <div className="flex flex-col">
                <p className="font-bold text-gray-12 text-xs">
                  {session?.user.name}
                </p>
                <p className="text-gray-11 text-xs">{session?.user.email}</p>
              </div>
            </div>
            <a
              href="https://docs.starliner.dev"
              rel="noreferrer"
              target="_blank"
              className="flex flex-row items-center gap-2 rounded-md p-2 text-xs hover:bg-gray-3"
            >
              <p>Documentation</p>
            </a>
            <button
              type="button"
              className="flex cursor-pointer flex-row items-center gap-2 rounded-md p-2 text-xs hover:bg-gray-3"
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

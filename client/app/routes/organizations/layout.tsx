import { Outlet, redirect } from "react-router";
import { auth } from "~/utils/auth/server";
import type { Route } from "./+types/layout";

export async function loader({ request }: Route.LoaderArgs) {
  const session = await auth.api.getSession({
    headers: request.headers,
  });

  if (!session) {
    const url = new URL(request.url);
    return redirect(`/login?redirectTo=${encodeURIComponent(url.pathname)}`);
  }
}

export default function Page() {
  return (
    <div className="flex min-h-screen">
      <div className="bg-mauve-4 wiggle-pattern w-1/2"></div>
      <div className="flex w-1/2 items-center justify-center p-16 shadow-md">
        <Outlet />
      </div>
    </div>
  );
}

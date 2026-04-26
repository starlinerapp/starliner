import type { Route } from "./+types/layout";
import { Outlet, redirect } from "react-router";
import React from "react";
import { auth } from "~/utils/auth/server";

export async function loader({ request }: Route.LoaderArgs) {
  const session = await auth.api.getSession({
    headers: request.headers,
  });

  if (session) {
    const url = new URL(request.url);
    const redirectTo = url.searchParams.get("redirectTo") || "/";
    return redirect(redirectTo);
  }
}

export default function Page() {
  return (
    <div className="flex min-h-dvh flex-col md:flex-row">
      <div className="bg-mauve-4 wiggle-pattern hidden md:block md:w-1/2" />
      <div className="flex w-full flex-1 items-center justify-center px-4 py-8 shadow-md sm:px-6 md:w-1/2 md:p-16">
        <Outlet />
      </div>
    </div>
  );
}

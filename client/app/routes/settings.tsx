import React from "react";
import type { Route } from "../../.react-router/types/app/+types/root";
import { auth } from "~/lib/auth.server";
import { redirect } from "react-router";
import Sidebar from "~/components/organisms/sidebar/Sidebar";

export function meta() {
  return [{ title: "Starliner" }, { name: "description", content: "" }];
}

export async function loader({ request }: Route.LoaderArgs) {
  const session = await auth.api.getSession({
    headers: request.headers,
  });

  if (!session) {
    return redirect("/login");
  }
}

export default function Settings() {
  return (
    <div>
      <Sidebar />
    </div>
  );
}

import React from "react";
import { redirect } from "react-router";
import Sidebar from "~/components/organisms/sidebar/Sidebar";
import type { Route } from "../../.react-router/types/app/+types/root";
import { auth } from "~/utils/auth/server";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";

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

export default function Home() {
  const trpc = useTRPC();
  const { data: root } = useQuery(trpc.root.getRoot.queryOptions());
  const { data: user } = useQuery(trpc.user.getUser.queryOptions());
  console.log(root);
  console.log(user);

  return (
    <div>
      <Sidebar />
    </div>
  );
}

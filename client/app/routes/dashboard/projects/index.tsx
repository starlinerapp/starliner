import React from "react";
import type { Route } from "../../../../.react-router/types/app/+types/root";
import { redirect } from "react-router";

export function meta() {
  return [{ title: "Starliner" }, { name: "description", content: "" }];
}

export async function loader({ params }: Route.LoaderArgs) {
  throw redirect(`/${params.slug}/projects/all`);
}

export default function Index() {
  return <></>;
}

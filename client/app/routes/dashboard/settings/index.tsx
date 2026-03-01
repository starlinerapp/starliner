import React from "react";
import type { Route } from "./+types/index";
import { redirect } from "react-router";

export async function loader(loaderArgs: Route.LoaderArgs) {
  throw redirect(`/${loaderArgs.params.slug}/settings/organization`);
}

export default function Settings() {
  return <></>;
}

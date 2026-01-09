import React from "react";
import type { Route } from "./+types/index";
import { redirect } from "react-router";

export async function loader({ params }: Route.LoaderArgs) {
  throw redirect(`/${params.slug}/projects`);
}

export default function Index() {
  return <></>;
}

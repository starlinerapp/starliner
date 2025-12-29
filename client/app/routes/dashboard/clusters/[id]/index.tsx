import React from "react";
import type { Route } from "../../../../../.react-router/types/app/+types/root";
import { redirect } from "react-router";

export async function loader({ params }: Route.LoaderArgs) {
  throw redirect(`/${params.slug}/clusters/${params.id}/general`);
}

export default function Cluster() {
  return <></>;
}

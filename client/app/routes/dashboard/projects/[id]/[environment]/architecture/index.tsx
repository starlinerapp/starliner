import React from "react";
import { redirect } from "react-router";
import { caller } from "~/utils/trpc/server";
import type { Route } from "./+types/layout";

export async function loader(loaderArgs: Route.LoaderArgs) {
  const { params } = loaderArgs;
  const trpc = await caller(loaderArgs);

  if (params.environment) {
    throw redirect(
      `/${params.slug}/projects/${params.id}/${params.environment}/architecture/git`,
    );
  }

  const project = await trpc.project.getProject({
    id: Number(params.id),
  });

  throw redirect(
    `/${params.slug}/projects/${params.id}/${project.environments[0].slug}/architecture/git`,
  );
}

export default function Architecture() {
  return <></>;
}

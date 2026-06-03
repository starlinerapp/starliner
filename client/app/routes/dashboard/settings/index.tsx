import { redirect } from "react-router";
import type { Route } from "./+types/index";

export async function loader(loaderArgs: Route.LoaderArgs) {
  throw redirect(`/${loaderArgs.params.slug}/settings/organization/members`);
}

export default function Settings() {
  return null;
}

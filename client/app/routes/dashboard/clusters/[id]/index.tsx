import { redirect } from "react-router";
import type { Route } from "./+types/index";

export async function loader({ params }: Route.LoaderArgs) {
  throw redirect(`/${params.slug}/clusters/${params.id}/general`);
}

export default function Cluster() {
  return null;
}

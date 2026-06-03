import { redirect } from "react-router";
import type { Route } from "./+types/index";

export async function loader({ params }: Route.LoaderArgs) {
  throw redirect(`/${params.slug}/clusters/all`);
}

export default function Index() {
  return null;
}

import { redirect } from "react-router";
import type { Route } from "./+types/index";

export function meta() {
  return [{ title: "Starliner" }, { name: "description", content: "" }];
}

export async function loader({ params }: Route.LoaderArgs) {
  throw redirect(`/${params.slug}/projects/all`);
}

export default function Index() {
  return null;
}

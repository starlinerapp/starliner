import type { Route } from "./+types/private-key";
import { getServerSession } from "~/utils/auth/server";
import { clusterApiFactory } from "~/server/api/client";

export const loader = async ({ params, request }: Route.LoaderArgs) => {
  const { id } = params;
  const session = await getServerSession(request);

  if (!session) {
    return new Response("Unauthorized", { status: 401 });
  }

  const response = await clusterApiFactory.getClusterPrivateKey(
    session?.user.id,
    Number(id),
  );

  return new Response(response.data, {
    status: 200,
    headers: {
      "Content-Type": "application/octet-stream",
      "Content-Disposition": `attachment; filename="private-key.pem"`,
    },
  });
};

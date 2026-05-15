import { serverEnv } from "~/env.server";

export type AuthUserRow = { id: string; name: string; email: string };

export async function fetchAuthUsersByIds(
  ids: string[],
): Promise<Map<string, AuthUserRow>> {
  if (ids.length === 0) {
    return new Map();
  }

  const res = await fetch(`${serverEnv.AUTH_URL}/internal/users`, {
    method: "POST",
    body: JSON.stringify({ ids }),
  });

  if (!res.ok) {
    console.error("failed to fetch users from auth service:", res.status);
    return new Map();
  }

  const body = (await res.json()) as { users?: AuthUserRow[] };
  const users = body.users ?? [];
  return new Map(users.map((u) => [u.id, u]));
}

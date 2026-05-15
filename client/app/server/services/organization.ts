import { fetchAuthUsersByIds } from "~/server/services/authUsers.server";

export async function enrichMembersWithAuthDetails<
  T extends { better_auth_id: string },
>(members: T[]) {
  const betterAuthIds = members.map((m) => m.better_auth_id);
  if (betterAuthIds.length === 0) return [];

  const authUserMap = await fetchAuthUsersByIds(betterAuthIds);

  return members.map((m) => ({
    ...m,
    name: authUserMap.get(m.better_auth_id)?.name ?? "",
    email: authUserMap.get(m.better_auth_id)?.email ?? "",
  }));
}

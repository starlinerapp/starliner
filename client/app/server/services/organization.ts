import { db } from "~/db";
import { user } from "~/db/schema";
import { inArray } from "drizzle-orm";

export async function enrichMembersWithAuthDetails<
  T extends { better_auth_id: string },
>(members: T[]) {
  const betterAuthIds = members.map((m) => m.better_auth_id);
  if (betterAuthIds.length === 0) return [];

  const authUsers = await db
    .select({ id: user.id, name: user.name, email: user.email })
    .from(user)
    .where(inArray(user.id, betterAuthIds));

  const authUserMap = new Map(authUsers.map((u) => [u.id, u]));

  return members.map((m) => ({
    ...m,
    name: authUserMap.get(m.better_auth_id)?.name ?? "",
    email: authUserMap.get(m.better_auth_id)?.email ?? "",
  }));
}

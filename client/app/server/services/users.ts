import { usersApiFactory } from "~/server/api/clients/auth";

export async function enrichMembersWithAuthDetails<
  T extends { better_auth_id: string },
>(
  members: T[],
): Promise<
  Array<
    T & {
      name: string;
      email: string;
    }
  >
> {
  if (members.length === 0) {
    return [];
  }

  try {
    const response = await usersApiFactory.bulkUserLookup({
      ids: members.map((m) => m.better_auth_id),
    });

    const users = response.data.users ?? [];

    const authUserMap = new Map(users.map((u) => [u.id, u]));

    return members.map((m) => ({
      ...m,
      name: authUserMap.get(m.better_auth_id)?.name ?? "",
      email: authUserMap.get(m.better_auth_id)?.email ?? "",
    }));
  } catch (error) {
    console.error("failed to fetch users from auth service:", error);

    return members.map((m) => ({
      ...m,
      name: "",
      email: "",
    }));
  }
}

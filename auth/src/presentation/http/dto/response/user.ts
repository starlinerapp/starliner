import type { UserProfile } from "~/domain/value/user";

export type BulkUserLookupResponse = {
  users: UserProfile[];
};

export function newBulkUserLookupResponse(
  users: UserProfile[],
): BulkUserLookupResponse {
  return { users };
}

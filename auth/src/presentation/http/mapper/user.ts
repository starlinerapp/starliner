import type { BulkUserLookupResponse } from "../dto/response/user";
import type { UserProfile } from "../../../domain/value/user";

export function toBulkUserLookupResponse(
  users: UserProfile[],
): BulkUserLookupResponse {
  return { users };
}

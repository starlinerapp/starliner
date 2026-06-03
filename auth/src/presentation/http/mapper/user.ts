import type { UserProfile } from "../../../domain/value/user";
import type { BulkUserLookupResponse } from "../dto/response/user";

export function toBulkUserLookupResponse(
  users: UserProfile[],
): BulkUserLookupResponse {
  return { users };
}

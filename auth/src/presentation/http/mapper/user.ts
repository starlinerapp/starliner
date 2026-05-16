import type { UserProfile } from "~/domain/value/user";
import type { BulkUserLookupResponse } from "~/presentation/http/dto/response/user";

export function toBulkUserLookupResponse(
  users: UserProfile[],
): BulkUserLookupResponse {
  return { users };
}

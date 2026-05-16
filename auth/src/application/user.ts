import { toUserProfiles, type UserProfile } from "~/domain/value/user";
import type { UserRepository } from "~/domain/repository/interface/user";

export const MAX_BULK_USER_LOOKUP = 200;

export class UserLookupError extends Error {
  constructor(readonly code: "too_many_ids") {
    super(code);
    this.name = "UserLookupError";
  }
}

export class UserApplication {
  constructor(private readonly userRepository: UserRepository) {}

  async getUsersByIds(ids: string[]): Promise<UserProfile[]> {
    const uniqueIds = [...new Set(ids)];

    if (uniqueIds.length === 0) {
      return [];
    }

    if (uniqueIds.length > MAX_BULK_USER_LOOKUP) {
      throw new UserLookupError("too_many_ids");
    }

    const users = await this.userRepository.findByIds(uniqueIds);
    return toUserProfiles(users);
  }
}

import type { User } from "../../entity/user";

export interface UserRepository {
  findByIds(ids: string[]): Promise<User[]>;
}

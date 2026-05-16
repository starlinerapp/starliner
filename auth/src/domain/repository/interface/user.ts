import type { User } from "~/domain/entity/user";

export interface UserRepository {
  findByIds(ids: string[]): Promise<User[]>;
}

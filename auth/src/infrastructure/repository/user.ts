import { inArray } from "drizzle-orm";
import type { db } from "../db";
import type { User } from "../../domain/entity/user";
import type { UserRepository } from "../../domain/repository/interface/user";
import { user as userTable } from "../../infrastructure/db/schema";

type Db = typeof db;

function toUser(row: {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
  image: string | null;
  createdAt: Date;
  updatedAt: Date;
}): User {
  return {
    id: row.id,
    name: row.name,
    email: row.email,
    emailVerified: row.emailVerified,
    image: row.image,
    createdAt: row.createdAt,
    updatedAt: row.updatedAt,
  };
}

export class DrizzleUserRepository implements UserRepository {
  constructor(private readonly db: Db) {}

  async findByIds(ids: string[]): Promise<User[]> {
    const rows = await this.db
      .select()
      .from(userTable)
      .where(inArray(userTable.id, ids));

    return rows.map(toUser);
  }
}

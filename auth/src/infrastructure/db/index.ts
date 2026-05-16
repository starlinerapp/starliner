import "dotenv/config";
import { drizzle } from "drizzle-orm/node-postgres";
import { serverEnv } from "../../env.server";
import { account, session, user, verification } from "./schema";

export const db = drizzle(serverEnv.AUTH_DATABASE_URL, {
  schema: { user, session, account, verification },
});

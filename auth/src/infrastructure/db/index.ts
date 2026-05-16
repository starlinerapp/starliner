import "dotenv/config";
import { drizzle } from "drizzle-orm/node-postgres";
import { serverEnv } from "~/env.server";

export const db = drizzle(serverEnv.AUTH_DATABASE_URL);

import "dotenv/config";
import { drizzle } from "drizzle-orm/node-postgres";
import { serverEnv } from "~/env.server";
import * as schema from "~/infrastructure/db/schema";

export const db = drizzle(serverEnv.AUTH_DATABASE_URL, { schema });

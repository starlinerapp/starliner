import "dotenv/config";
import { drizzle } from "drizzle-orm/node-postgres";
import { serverEnv } from "../../env.server";
import { schema } from "better-auth/client/plugins";

export const db = drizzle(serverEnv.AUTH_DATABASE_URL, { schema });

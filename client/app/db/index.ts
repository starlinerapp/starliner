import "dotenv/config";
import { drizzle } from "drizzle-orm/node-postgres";
// @ts-expect-error .ts extension required by custom node server
import { serverEnv } from "../env.server.ts";

export const db = drizzle(serverEnv.AUTH_DATABASE_URL);

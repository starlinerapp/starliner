import "dotenv/config";
import { drizzle } from "drizzle-orm/node-postgres";
// @ts-expect-error .ts required by server.js
import { serverEnv } from "../env.server.ts";

export const db = drizzle(serverEnv.AUTH_DATABASE_URL);

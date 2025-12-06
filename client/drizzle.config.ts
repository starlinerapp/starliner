import "dotenv/config";
import { defineConfig } from "drizzle-kit";
import { serverEnv } from "~/env.server";

export default defineConfig({
  out: "./drizzle",
  schema: "./app/db/schema.ts",
  dialect: "postgresql",
  dbCredentials: {
    url: serverEnv.AUTH_DATABASE_URL,
  },
});

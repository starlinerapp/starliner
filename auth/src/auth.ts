import { betterAuth } from "better-auth";
import { bearer } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "~/db";
import * as schema from "~/db/schema";

export const auth = betterAuth({
  baseURL: "https://auth.dev.starliner.app/api/auth",
  trustedOrigins: ["https://dev.starliner.app"],
  database: drizzleAdapter(db, {
    provider: "pg",
    schema,
  }),
  emailAndPassword: {
    enabled: true,
    requireEmailVerification: true,
    revokeSessionsOnPasswordReset: true,
    sendResetPassword: async ({ user, url }) => {},
  },
  emailVerification: {
    sendVerificationEmail: async ({ user, url }) => {},
  },
  plugins: [bearer()],
});

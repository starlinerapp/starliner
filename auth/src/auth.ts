import { betterAuth } from "better-auth";
import { bearer } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "~/db";
import * as schema from "~/db/schema";

export const auth = betterAuth({
  trustedOrigins: [
    "https://dev.starliner.app",
    "https://auth.dev.starliner.app",
  ],
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
  session: {
    cookieCache: {
      enabled: true,
      maxAge: 5 * 60,
    },
  },
  advanced: {
    cookiePrefix: "starliner",
    crossSubDomainCookies: {
      enabled: true,
      domain: "dev.starliner.app",
    },
    defaultCookieAttributes: {
      sameSite: "none",
      secure: true,
      httpOnly: true,
    },
  },
  plugins: [bearer()],
});

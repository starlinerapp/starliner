import { betterAuth } from "better-auth";
import { bearer } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "~/db";
import * as schema from "~/db/schema";
import { serverEnv } from "~/env.server";

const clientOrigin = new URL(serverEnv.CLIENT_BASE_URL).origin;
const authOrigin = new URL(serverEnv.AUTH_PUBLIC_URL).origin;

export const auth = betterAuth({
  trustedOrigins: [clientOrigin, authOrigin],
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
      domain: new URL(serverEnv.CLIENT_BASE_URL).hostname,
    },
    defaultCookieAttributes: {
      sameSite: "none",
      secure: true,
      httpOnly: true,
    },
  },
  plugins: [bearer()],
});

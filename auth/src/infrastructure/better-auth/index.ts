import { betterAuth } from "better-auth";
import { bearer } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import type { AuthService } from "~/domain/port/auth";
import { serverEnv } from "~/env.server";
import type { db } from "~/infrastructure/db";
import * as schema from "~/infrastructure/db/schema";

type Db = typeof db;

export function createBetterAuth(db: Db): AuthService {
  const clientOrigin = new URL(serverEnv.CLIENT_BASE_URL).origin;
  const authOrigin = new URL(serverEnv.AUTH_PUBLIC_URL).origin;

  return betterAuth({
    trustedOrigins: [clientOrigin, authOrigin],
    database: drizzleAdapter(db, {
      provider: "pg",
      schema,
    }),
    emailAndPassword: {
      enabled: true,
      requireEmailVerification: true,
      revokeSessionsOnPasswordReset: true,
      sendResetPassword: async () => {},
    },
    emailVerification: {
      sendVerificationEmail: async () => {},
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
}

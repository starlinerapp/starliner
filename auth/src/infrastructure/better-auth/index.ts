import { betterAuth } from "better-auth";
import { bearer } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import type { EmailApplication } from "src/application/email";
import type { db } from "../db";
import type { AuthService } from "../../domain/port/auth";
import { serverEnv } from "../../env.server";
import { schema } from "better-auth/client/plugins";

type Db = typeof db;

export function createBetterAuth(
  db: Db,
  emailApplication: EmailApplication,
): AuthService {
  const clientOrigin = new URL(serverEnv.CLIENT_BASE_URL).origin;
  const authOrigin = new URL(serverEnv.AUTH_PUBLIC_URL).origin;

  const authBaseUrl = new URL(
    "/api/auth",
    serverEnv.AUTH_PUBLIC_URL,
  ).toString();

  return betterAuth({
    baseURL: authBaseUrl,
    trustedOrigins: [clientOrigin, authOrigin],
    database: drizzleAdapter(db, {
      provider: "pg",
      schema,
    }),
    emailAndPassword: {
      enabled: true,
      requireEmailVerification: true,
      revokeSessionsOnPasswordReset: true,
      sendResetPassword: async ({ user, url }) => {
        await emailApplication.sendResetPassword(user.email, url);
      },
    },
    emailVerification: {
      autoSignInAfterVerification: true,
      sendVerificationEmail: async ({ user, url }) => {
        await emailApplication.sendVerificationEmail(user.email, url);
      },
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

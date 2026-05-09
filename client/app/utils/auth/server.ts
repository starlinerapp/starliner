import { betterAuth } from "better-auth";
import { bearer, openAPI } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { internalApiFactory } from "~/server/api/client";
import { db } from "~/db";
import * as schema from "~/db/schema";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    schema,
  }),
  emailAndPassword: {
    enabled: true,
    requireEmailVerification: true,
    revokeSessionsOnPasswordReset: true,
    sendResetPassword: async ({ user, url }) => {
      await internalApiFactory.sendResetPassword(user.id, {
        to: user.email,
        resetUrl: url,
      });
    },
  },
  emailVerification: {
    sendVerificationEmail: async ({ user, url }) => {
      await internalApiFactory.sendVerificationEmail(user.id, {
        to: user.email,
        verificationUrl: url,
      });
    },
  },
  plugins: [
    bearer(),
    ...(process.env.NODE_ENV !== "production" ? [openAPI()] : []),
  ],
});

export const getServerSession = async (request: Request) => {
  return await auth.api.getSession({
    headers: request.headers,
  });
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
type Session = typeof auth.$Infer.Session;

import { betterAuth } from "better-auth";
import { bearer } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "~/db";
import * as schema from "~/db/schema";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    schema,
  }),
  emailAndPassword: {
    enabled: true,
  },
  plugins: [bearer()],
});

export const getServerSession = async (request: Request) => {
  return await auth.api.getSession({
    headers: request.headers,
  });
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
type Session = typeof auth.$Infer.Session;

import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
// @ts-expect-error .ts required by server.js
import { db } from "../../db/index.ts";
// @ts-expect-error .ts required by server.js
import * as schema from "../../db/schema.ts";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    schema,
  }),
  emailAndPassword: {
    enabled: true,
  },
});

export const getServerSession = async (request: Request) => {
  return await auth.api.getSession({
    headers: request.headers,
  });
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
type Session = typeof auth.$Infer.Session;

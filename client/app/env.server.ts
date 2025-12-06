import { z } from "zod";

const envSchema = z.object({
  BETTER_AUTH_SECRET: z.string(),
  BETTER_AUTH_URL: z.string(),
  AUTH_DATABASE_URL: z.string(),
  SERVER_BASE_URL: z.string(),
  SERVER_BASIC_AUTH_USER: z.string(),
  SERVER_BASIC_AUTH_PASSWORD: z.string(),
});

function createServerEnv() {
  const parsed = envSchema.safeParse(process.env);

  if (!parsed.success) {
    console.error("Invalid environment variables:", parsed.error.issues);
    throw new Error("Invalid environment variables");
  }

  return parsed.data;
}

export const serverEnv = createServerEnv();
export type ServerEnv = z.infer<typeof envSchema>;

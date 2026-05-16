import { UserApplication } from "~/application/user";
import { db } from "~/infrastructure/db";
import { createBetterAuth } from "~/infrastructure/better-auth";
import { DrizzleUserRepository } from "~/infrastructure/repository/user";
import { createApp } from "~/presentation/http/app";

export function bootstrap() {
  const userRepository = new DrizzleUserRepository(db);
  const userApplication = new UserApplication(userRepository);
  const auth = createBetterAuth(db);

  const app = createApp({ userApplication, auth });

  return { app };
}

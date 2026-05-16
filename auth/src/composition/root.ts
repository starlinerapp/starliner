import { EmailApplication } from "~/application/email";
import { UserApplication } from "~/application/user";
import { db } from "~/infrastructure/db";
import { createAuthApiService } from "~/infrastructure/api/client";
import { createBetterAuth } from "~/infrastructure/better-auth";
import { DrizzleUserRepository } from "~/infrastructure/repository/user";
import { createApp } from "~/presentation/http/app";

export function bootstrap() {
  const userRepository = new DrizzleUserRepository(db);
  const userApplication = new UserApplication(userRepository);

  const apiService = createAuthApiService();
  const emailApplication = new EmailApplication(apiService);
  const auth = createBetterAuth(db, emailApplication);

  const app = createApp({ userApplication, auth });

  return { app };
}

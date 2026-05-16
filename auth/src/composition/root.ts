import { DrizzleUserRepository } from "../infrastructure/repository/user";
import { db } from "../infrastructure/db";
import { UserApplication } from "../application/user";
import { createAuthApiService } from "../infrastructure/api/client";
import { EmailApplication } from "../application/email";
import { createBetterAuth } from "../infrastructure/better-auth";
import { createApp } from "../presentation/http/app";

export function bootstrap() {
  const userRepository = new DrizzleUserRepository(db);
  const userApplication = new UserApplication(userRepository);

  const apiService = createAuthApiService();
  const emailApplication = new EmailApplication(apiService);
  const auth = createBetterAuth(db, emailApplication);

  const app = createApp({ userApplication, auth });

  return { app };
}

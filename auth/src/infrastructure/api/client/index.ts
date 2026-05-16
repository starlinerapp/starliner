import {
  AuthApiFactory,
  Configuration,
} from "~/infrastructure/api/client/generated";
import { serverEnv } from "~/env.server";
import type { ApiService } from "~/domain/port/api";
import { createServerApiAxios } from "~/infrastructure/api/client/axios";

export class AuthApiService implements ApiService {
  constructor(
    private readonly authApiFactory: ReturnType<typeof AuthApiFactory>,
  ) {}

  async sendResetPasswordEmail(resetUrl: string, to: string): Promise<void> {
    await this.authApiFactory.sendResetPassword({
      resetUrl,
      to,
    });
  }

  async sendVerificationEmail(
    verificationUrl: string,
    to: string,
  ): Promise<void> {
    await this.authApiFactory.sendVerificationEmail({
      verificationUrl,
      to,
    });
  }
}

export function createAuthApiService(): AuthApiService {
  const configuration = new Configuration({
    basePath: `http://${serverEnv.SERVER_BASE_URL}`,
  });

  const authApiFactory = AuthApiFactory(
    configuration,
    undefined,
    createServerApiAxios(),
  );

  return new AuthApiService(authApiFactory);
}

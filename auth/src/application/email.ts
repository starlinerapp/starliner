import type { ApiService } from "~/domain/port/api";

export class EmailApplication {
  constructor(private readonly api: ApiService) {}

  async sendResetPassword(to: string, resetUrl: string): Promise<void> {
    await this.api.sendResetPasswordEmail(resetUrl, to);
  }

  async sendVerificationEmail(
    to: string,
    verificationUrl: string,
  ): Promise<void> {
    await this.api.sendVerificationEmail(verificationUrl, to);
  }
}

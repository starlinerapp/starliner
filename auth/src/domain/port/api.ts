export interface ApiService {
  sendResetPasswordEmail(resetUrl: string, to: string): Promise<void>;
  sendVerificationEmail(verificationUrl: string, to: string): Promise<void>;
}

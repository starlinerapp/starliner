export {};

declare global {
  interface Window {
    ENV: {
      SENTRY_DSN_CLIENT: string;
      ENVIRONMENT: string;
      AUTH_PUBLIC_URL: string;
    };
  }
}

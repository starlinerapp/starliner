import { createAuthClient } from "better-auth/react";

type AuthClient = ReturnType<typeof createAuthClient>;

let authClient: AuthClient | undefined;

export function getAuthClient(): AuthClient {
  if (!authClient) {
    const url =
      typeof window !== "undefined"
        ? window.ENV.AUTH_PUBLIC_URL
        : process.env.AUTH_PUBLIC_URL;

    if (!url) {
      throw new Error("AUTH_PUBLIC_URL is not set");
    }

    authClient = createAuthClient({
      baseURL: `${url.replace(/\/$/, "")}/api/auth`,
    });
  }
  return authClient;
}

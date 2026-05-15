import React from "react";
import { createContext, useContext, useMemo, type ReactNode } from "react";
import { createAuthClient } from "better-auth/react";

type AuthClient = ReturnType<typeof createAuthClient>;

const AuthClientContext = createContext<AuthClient | null>(null);

export function AuthClientProvider({
  authPublicOrigin,
  children,
}: {
  authPublicOrigin: string;
  children: ReactNode;
}) {
  const client = useMemo(() => {
    const origin = authPublicOrigin.replace(/\/$/, "");
    return createAuthClient({
      baseURL: `${origin}/api/auth`,
    });
  }, [authPublicOrigin]);

  return (
    <AuthClientContext.Provider value={client}>
      {children}
    </AuthClientContext.Provider>
  );
}

export function useAuthClient(): AuthClient {
  const client = useContext(AuthClientContext);
  if (!client) {
    throw new Error("useAuthClient must be used within AuthClientProvider");
  }
  return client;
}

import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

type CreateAuthClientOptions = {
  baseURL: string;
  sessionOptions?: {
    refetchInterval?: number;
  };
};

const mockCreateAuthClient = vi.fn<
  (options: CreateAuthClientOptions) => {
    useSession: ReturnType<typeof vi.fn>;
    signIn: { email: ReturnType<typeof vi.fn> };
    signOut: ReturnType<typeof vi.fn>;
  }
>(() => ({
  useSession: vi.fn(),
  signIn: { email: vi.fn() },
  signOut: vi.fn(),
}));

vi.mock("better-auth/react", () => ({
  createAuthClient: mockCreateAuthClient,
}));

const AUTH_PUBLIC_URL = "https://auth.test.starliner.app/";

function stubBrowserEnv(url = AUTH_PUBLIC_URL) {
  vi.stubGlobal("window", {
    ENV: {
      SENTRY_DSN_CLIENT: "",
      ENVIRONMENT: "test",
      AUTH_PUBLIC_URL: url,
    },
  } as Window & typeof globalThis);
}

async function importAuthClient() {
  return import("~/utils/auth/client");
}

describe("getAuthClient", () => {
  beforeEach(() => {
    vi.resetModules();
    mockCreateAuthClient.mockClear();
    delete process.env.AUTH_PUBLIC_URL;
    stubBrowserEnv();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("refetches the session before the server cookie cache expires", async () => {
    const { getAuthClient } = await importAuthClient();

    getAuthClient();

    const config = mockCreateAuthClient.mock.calls[0]?.[0];
    // Auth service cookie cache maxAge is 5 minutes; client refreshes at 4 minutes.
    expect(config?.sessionOptions?.refetchInterval).toBe(240);
  });

  it("returns a singleton auth client instance", async () => {
    const { getAuthClient } = await importAuthClient();
    const first = getAuthClient();
    const second = getAuthClient();

    expect(first).toBe(second);
    expect(mockCreateAuthClient).toHaveBeenCalledTimes(1);
  });

  it("throws when AUTH_PUBLIC_URL is not set in the browser", async () => {
    vi.stubGlobal("window", {
      ENV: {
        SENTRY_DSN_CLIENT: "",
        ENVIRONMENT: "test",
        AUTH_PUBLIC_URL: "",
      },
    } as Window & typeof globalThis);

    const { getAuthClient } = await importAuthClient();

    expect(() => getAuthClient()).toThrow("AUTH_PUBLIC_URL is not set");
  });
});

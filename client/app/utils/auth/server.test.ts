import type { IncomingHttpHeaders } from "node:http";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const mockGetSession = vi.fn();
const mockCreateAuthClient = vi.fn(() => ({
  getSession: mockGetSession,
}));

vi.mock("~/env.server", () => ({
  serverEnv: {
    AUTH_URL: "https://auth.test.starliner.app/",
    AUTH_PUBLIC_URL: "https://auth.test.starliner.app",
    SERVER_BASE_URL: "http://localhost:8080",
    SERVER_BASIC_AUTH_USER: "test-user",
    SERVER_BASIC_AUTH_PASSWORD: "test-pass",
  },
}));

vi.mock("better-auth/client", () => ({
  createAuthClient: mockCreateAuthClient,
}));

const {
  toSessionLookupHeaders,
  getSessionFromHeaders,
  getSessionFromNodeHeaders,
  getServerSession,
} = await import("~/utils/auth/server");

const sessionUser = {
  id: "user-1",
  email: "user@example.com",
  name: "Test User",
  emailVerified: true,
  createdAt: new Date(),
  updatedAt: new Date(),
};

const sessionRecord = {
  id: "session-1",
  userId: "user-1",
  expiresAt: new Date(Date.now() + 60_000),
  token: "session-token",
  createdAt: new Date(),
  updatedAt: new Date(),
};

describe("auth server utilities", () => {
  beforeEach(() => {
    mockGetSession.mockReset();
    vi.spyOn(console, "error").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("toSessionLookupHeaders", () => {
    it("copies forwardable headers and drops hop-by-hop headers", () => {
      const headers: IncomingHttpHeaders = {
        cookie: "starliner.session_token=abc",
        host: "client.test.starliner.app",
        connection: "Upgrade",
        upgrade: "websocket",
        "sec-websocket-key": "dGhlIHNhbXBsZSBub25jZQ==",
        "sec-websocket-version": "13",
      };

      const result = toSessionLookupHeaders(headers);

      expect(result.get("cookie")).toBe("starliner.session_token=abc");
      expect(result.get("host")).toBe("client.test.starliner.app");
      expect(result.has("connection")).toBe(false);
      expect(result.has("upgrade")).toBe(false);
      expect(result.has("sec-websocket-key")).toBe(false);
      expect(result.has("sec-websocket-version")).toBe(false);
    });

    it("appends array header values", () => {
      const result = toSessionLookupHeaders({
        "set-cookie": ["a=1", "b=2"],
      });

      expect(result.getSetCookie()).toEqual(["a=1", "b=2"]);
    });

    it("skips undefined header values", () => {
      const result = toSessionLookupHeaders({
        cookie: undefined,
        authorization: "Bearer token",
      });

      expect(result.has("cookie")).toBe(false);
      expect(result.get("authorization")).toBe("Bearer token");
    });
  });

  describe("getSessionFromHeaders", () => {
    it("forwards request headers to Better Auth getSession", async () => {
      const headers = new Headers({
        cookie: "starliner.session_token=initial",
      });
      mockGetSession.mockResolvedValue({
        data: { user: sessionUser, session: sessionRecord },
        error: null,
      });

      await getSessionFromHeaders(headers);

      expect(mockGetSession).toHaveBeenCalledWith({
        fetchOptions: { headers },
      });
    });

    it("returns user and session when Better Auth responds successfully", async () => {
      mockGetSession.mockResolvedValue({
        data: { user: sessionUser, session: sessionRecord },
        error: null,
      });

      const result = await getSessionFromHeaders(new Headers());

      expect(result).toEqual({
        user: sessionUser,
        session: sessionRecord,
      });
    });

    it("returns null and logs when Better Auth reports an error", async () => {
      mockGetSession.mockResolvedValue({
        data: null,
        error: { message: "invalid session" },
      });

      const result = await getSessionFromHeaders(new Headers());

      expect(result).toBeNull();
      expect(console.error).toHaveBeenCalledWith(
        "getSessionFromHeaders failed:",
        { message: "invalid session" },
      );
    });

    it("returns null when the session payload is incomplete", async () => {
      mockGetSession.mockResolvedValue({
        data: { user: sessionUser, session: null },
        error: null,
      });

      await expect(getSessionFromHeaders(new Headers())).resolves.toBeNull();
    });
  });

  describe("getSessionFromNodeHeaders", () => {
    it("sanitizes Node headers before session lookup", async () => {
      mockGetSession.mockResolvedValue({
        data: { user: sessionUser, session: sessionRecord },
        error: null,
      });

      await getSessionFromNodeHeaders({
        cookie: "starliner.session_token=test",
        upgrade: "websocket",
      });

      const forwardedHeaders = mockGetSession.mock.calls[0]?.[0]?.fetchOptions
        ?.headers as Headers;
      expect(forwardedHeaders.get("cookie")).toBe(
        "starliner.session_token=test",
      );
      expect(forwardedHeaders.has("upgrade")).toBe(false);
    });
  });

  describe("getServerSession", () => {
    it("reads session cookies from a fetch request", async () => {
      mockGetSession.mockResolvedValue({
        data: { user: sessionUser, session: sessionRecord },
        error: null,
      });

      const request = new Request(
        "https://client.test.starliner.app/dashboard",
        {
          headers: {
            cookie: "starliner.session_token=from-request",
          },
        },
      );

      const result = await getServerSession(request);

      expect(result?.session.token).toBe("session-token");
      const forwardedHeaders = mockGetSession.mock.calls[0]?.[0]?.fetchOptions
        ?.headers as Headers;
      expect(forwardedHeaders.get("cookie")).toBe(
        "starliner.session_token=from-request",
      );
    });
  });
});

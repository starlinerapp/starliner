import { createAuthClient } from "better-auth/client";
import type { IncomingHttpHeaders } from "node:http";
import { serverEnv } from "~/env.server";

const authBaseURL = `${serverEnv.AUTH_URL.replace(/\/$/, "")}/api/auth`;

const nodeAuthClient = createAuthClient({
  baseURL: authBaseURL,
});

/** Headers that must not be forwarded to a normal HTTP fetch (e.g. from a WS upgrade). */
const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "upgrade",
  "keep-alive",
  "transfer-encoding",
  "te",
  "trailer",
  "proxy-connection",
  "sec-websocket-key",
  "sec-websocket-version",
  "sec-websocket-extensions",
  "sec-websocket-protocol",
  "sec-websocket-accept",
]);

export function toSessionLookupHeaders(headers: IncomingHttpHeaders): Headers {
  const sessionHeaders = new Headers();

  for (const [key, value] of Object.entries(headers)) {
    if (value === undefined || HOP_BY_HOP_HEADERS.has(key.toLowerCase())) {
      continue;
    }

    if (Array.isArray(value)) {
      for (const entry of value) {
        sessionHeaders.append(key, entry);
      }
    } else {
      sessionHeaders.set(key, value);
    }
  }

  return sessionHeaders;
}

export async function getSessionFromNodeHeaders(headers: IncomingHttpHeaders) {
  return getSessionFromHeaders(toSessionLookupHeaders(headers));
}

export async function getSessionFromHeaders(headers: Headers) {
  const result = await nodeAuthClient.getSession({
    fetchOptions: { headers },
  });

  if (result.error) {
    return null;
  }

  const data = result.data;
  if (!data?.user || !data.session) {
    return null;
  }

  return { user: data.user, session: data.session };
}

export async function getServerSession(request: Request) {
  return getSessionFromHeaders(request.headers);
}

export const auth = {
  api: {
    getSession: (opts: { headers: Headers }) =>
      getSessionFromHeaders(opts.headers),
  },
};

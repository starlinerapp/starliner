import { createAuthClient } from "better-auth/client";
import { serverEnv } from "~/env.server";

const authBaseURL = `${serverEnv.AUTH_URL.replace(/\/$/, "")}/api/auth`;

const nodeAuthClient = createAuthClient({
  baseURL: authBaseURL,
});

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

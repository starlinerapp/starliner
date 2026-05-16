import { Hono } from "hono";
import { cors } from "hono/cors";
import type { UserApplication } from "~/application/user";
import type { AuthService } from "~/domain/port/auth";
import { serverEnv } from "~/env.server";
import { AuthHandler } from "~/presentation/http/handler/auth";
import { UserHandler } from "~/presentation/http/handler/user";

type AppDependencies = {
  userApplication: UserApplication;
  auth: AuthService;
};

export function createApp({ userApplication, auth }: AppDependencies) {
  const app = new Hono();
  const userHandler = new UserHandler(userApplication);
  const authHandler = new AuthHandler(auth);

  const internal = new Hono();
  internal.post("/users", userHandler.bulkLookup);
  app.route("/internal", internal);

  app.post("/users", userHandler.bulkLookup);

  app.use(
    "/api/auth/*",
    cors({
      origin: serverEnv.CLIENT_BASE_URL,
      allowHeaders: ["Content-Type", "Authorization"],
      allowMethods: ["POST", "GET", "OPTIONS"],
      exposeHeaders: ["Content-Length"],
      maxAge: 600,
      credentials: true,
    }),
  );

  app.on(["POST", "GET"], "/api/auth/*", authHandler.handle);

  return app;
}

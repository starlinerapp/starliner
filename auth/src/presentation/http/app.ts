import { swaggerUI } from "@hono/swagger-ui";
import { OpenAPIHono } from "@hono/zod-openapi";
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
  const app = new OpenAPIHono({
    defaultHook: (result, c) => {
      if (!result.success) {
        return c.json({ error: "invalid json" }, 400);
      }
    },
  });

  UserHandler.register(app, userApplication);

  app.doc("/docs", {
    openapi: "3.0.0",
    info: {
      title: "Starliner Auth API",
      version: "1.0.0",
    },
  });

  app.get("/ui", swaggerUI({ url: "/docs" }));

  const authHandler = new AuthHandler(auth);

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

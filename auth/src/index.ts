import { Hono } from "hono";
import { serve } from "@hono/node-server";
import { cors } from "hono/cors";
import { inArray } from "drizzle-orm";
import { auth } from "./auth";
import { db } from "~/db";
import { user } from "~/db/schema";
import { serverEnv } from "~/env.server";

const app = new Hono();

app.post("/users", async (c) => {
  let body: unknown;
  try {
    body = await c.req.json();
  } catch {
    return c.json({ error: "invalid json" }, 400);
  }

  const rawIds = (body as { ids?: unknown }).ids;
  const ids = Array.isArray(rawIds)
    ? rawIds.filter((x): x is string => typeof x === "string")
    : [];

  if (ids.length > 200) {
    return c.json({ error: "too many ids" }, 400);
  }
  if (ids.length === 0) {
    return c.json({ users: [] });
  }

  const rows = await db
    .select({ id: user.id, name: user.name, email: user.email })
    .from(user)
    .where(inArray(user.id, ids));

  return c.json({ users: rows });
});

app.route("/internal", app);

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

app.on(["POST", "GET"], "/api/auth/*", (c) => {
  return auth.handler(c.req.raw);
});

serve(app);

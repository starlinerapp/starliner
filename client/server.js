import compression from "compression";
import express from "express";
import morgan from "morgan";
import httpProxy from "http-proxy";
import * as http from "node:http";
import { auth } from "./app/utils/auth/server.ts";
import { fromNodeHeaders } from "better-auth/node";

const BUILD_PATH = "./build/server/index.js";
const DEVELOPMENT = process.env.NODE_ENV === "development";
const PORT = Number.parseInt(process.env.PORT || "5173");
const SERVER_BASE_URL = process.env.SERVER_BASE_URL;
const SERVER_BASIC_AUTH_USER = process.env.SERVER_BASIC_AUTH_USER;
const SERVER_BASIC_AUTH_PASSWORD = process.env.SERVER_BASIC_AUTH_PASSWORD;

const app = express();

app.use(compression());
app.disable("x-powered-by");

if (DEVELOPMENT) {
  console.log("Starting development server");
  const viteDevServer = await import("vite").then((vite) =>
    vite.createServer({
      server: { middlewareMode: true },
    }),
  );
  app.use(viteDevServer.middlewares);
  app.use(async (req, res, next) => {
    if (req.url?.startsWith("/ws")) {
      return next();
    }

    try {
      const source = await viteDevServer.ssrLoadModule("./server/app.ts");
      return await source.app(req, res, next);
    } catch (error) {
      if (typeof error === "object" && error instanceof Error) {
        viteDevServer.ssrFixStacktrace(error);
      }
      next(error);
    }
  });
} else {
  console.log("Starting production server");
  app.use(
    "/assets",
    express.static("build/client/assets", { immutable: true, maxAge: "1y" }),
  );
  app.use(morgan("tiny"));
  app.use(express.static("build/client", { maxAge: "1h" }));

  app.use(await import(BUILD_PATH).then((mod) => mod.app));
}

const wsProxy = httpProxy.createProxyServer({
  target: `ws://${SERVER_BASE_URL}`,
  ws: true,
  changeOrigin: true,
});

wsProxy.on("error", (err, req, socket) => {
  console.error("WebSocket proxy error:", err);

  if ("write" in socket) {
    socket.write("HTTP/1.1 502 Bad Gateway\r\n\r\n");
    socket.destroy();
  }
});

const server = http.createServer(app);
server.on("upgrade", async (req, socket, head) => {
  console.log("Proxying WebSocket request:", req.url);

  const session = await auth.api.getSession({
    headers: fromNodeHeaders(req.headers),
  });

  if (session.user.id) {
    req.headers["X-User-Id"] = session.user.id;
  }

  const credentials = Buffer.from(
    `${SERVER_BASIC_AUTH_USER}:${SERVER_BASIC_AUTH_PASSWORD}`,
  ).toString("base64");

  req.headers["authorization"] = `Basic ${credentials}`;

  wsProxy.ws(req, socket, head);
});

server.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});

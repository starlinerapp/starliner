import { reactRouter } from "@react-router/dev/vite";
import { sentryReactRouter } from "@sentry/react-router";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig((config) => ({
  plugins: [
    tailwindcss(),
    reactRouter(),
    tsconfigPaths(),
    sentryReactRouter(
      {
        org: "starliner",
        project: "client",
        authToken: process.env.SENTRY_AUTH_TOKEN_CLIENT,
      },
      config,
    ),
  ],

  server: {
    allowedHosts: [process.env.DOMAIN ?? "dev.starliner.app", "client"],
  },

  optimizeDeps: {
    exclude: ["@sentry/react-router"],
  },
}));

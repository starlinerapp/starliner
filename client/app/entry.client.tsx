import * as Sentry from "@sentry/react-router";
import { StrictMode, startTransition } from "react";
import { hydrateRoot } from "react-dom/client";
import { HydratedRouter } from "react-router/dom";

const BOT_PROBE_PATTERNS = ['.php', '/.git', '/wp-admin', '/wp-login', '/.env', '/.action'];

Sentry.init({
  dsn: window.ENV.SENTRY_DSN_CLIENT,
  environment: window.ENV.ENVIRONMENT,
  sendDefaultPii: true,
  beforeSend(event) {
    const url = event.request?.url ?? '';
    if (BOT_PROBE_PATTERNS.some((p) => url.includes(p))) return null;
    return event;
  },
  integrations: [
    Sentry.reactRouterTracingIntegration(),
    Sentry.replayIntegration(),
  ],
  enableLogs: true,
  tracesSampleRate: 1.0,
  tracePropagationTargets: [/^\//, /^https:\/\/yourserver\.io\/api/],
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0,
});

startTransition(() => {
  hydrateRoot(
    document,
    <StrictMode>
      <HydratedRouter />
    </StrictMode>,
  );
});

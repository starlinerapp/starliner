import * as Sentry from "@sentry/react-router";
import { StrictMode, startTransition } from "react";
import { hydrateRoot } from "react-dom/client";
import { isRouteErrorResponse } from "react-router";
import { HydratedRouter } from "react-router/dom";

Sentry.init({
  dsn: window.ENV.SENTRY_DSN_CLIENT,
  environment: window.ENV.ENVIRONMENT,
  sendDefaultPii: true,
  beforeSend(event, hint) {
    const error = hint.originalException;
    if (isRouteErrorResponse(error) && error.status === 404) return null;
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

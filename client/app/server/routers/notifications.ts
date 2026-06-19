import { randomUUID } from "node:crypto";
import { z } from "zod";
import { notificationsApiFactory } from "~/server/api/clients/server";
import { cache } from "~/server/services/cache";
import { streamSseJson } from "~/server/services/sse";
import { protectedProcedure } from "~/server/trpc";

export const notificationsRouter = {
  streamGlobalNotifications: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .subscription(async function* ({ input, ctx, signal }) {
      const userId = ctx.user?.id;

      const correlationCacheKey = `user:${userId}`;
      let correlationId = await cache.get(correlationCacheKey);

      if (!correlationId) {
        correlationId = randomUUID();
        await cache.set(correlationCacheKey, correlationId, 60 * 60 * 60);
      }

      yield* streamSseJson(
        () =>
          // @ts-expect-error OpenAPI doesn't support SSE
          notificationsApiFactory.streamGlobalNotifications(
            userId,
            input.organizationId,
            {
              responseType: "stream",
              signal,
              headers: { "X-Correlation-ID": correlationId },
            },
          ),
        signal,
      );
    }),
};

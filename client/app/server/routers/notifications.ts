import { z } from "zod";
import { notificationsApiFactory } from "~/server/api/clients/server";
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

      yield* streamSseJson(
        () =>
          // @ts-expect-error OpenAPI doesn't support SSE
          notificationsApiFactory.streamGlobalNotifications(
            userId,
            input.organizationId,
            { responseType: "stream", signal },
          ),
        signal,
      );
    }),
};

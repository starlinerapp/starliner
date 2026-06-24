import { z } from "zod";
import { buildApiFactory } from "~/server/api/clients/server";
import { streamSse } from "~/server/services/sse";
import { protectedProcedure } from "~/server/trpc";

export const buildRouter = {
  getBuildLogs: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await buildApiFactory
        .getBuildLogs(userId, input.id)
        .then((res) => res.data);
    }),
  streamBuildLogs: protectedProcedure
    .input(
      z.object({
        buildId: z.number(),
      }),
    )
    .subscription(async function* ({ input, ctx, signal }) {
      const userId = ctx.user?.id;

      yield* streamSse(
        () =>
          // @ts-expect-error OpenAPI doesn't support SSE
          buildApiFactory.streamBuildLogs(userId, input.buildId, {
            responseType: "stream",
            signal,
          }),
        signal,
        { trim: true },
      );
    }),
};

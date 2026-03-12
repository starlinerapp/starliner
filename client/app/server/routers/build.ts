import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { buildApiFactory } from "~/server/api/client";

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
};

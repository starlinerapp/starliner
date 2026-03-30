import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { githubFactory } from "~/server/api/client";

export const githubRouter = {
  getRepositories: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await githubFactory
        .getRepositories(userId, input.organizationId)
        .then((res) => res.data);
    }),
};

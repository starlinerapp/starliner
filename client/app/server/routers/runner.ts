import { z } from "zod";
import { runnerApiFactory } from "~/server/api/clients/server";
import { protectedProcedure } from "~/server/trpc";

export const runnerRouter = {
  createRunner: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await runnerApiFactory
        .createRunner(userId, input.organizationId)
        .then((res) => res.data);
    }),
  getOrganizationRunners: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await runnerApiFactory
        .getOrganizationRunners(userId, input.organizationId)
        .then((res) => res.data);
    }),
};

import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { clusterApiFactory } from "~/server/api/client";

export const clusterRouter = {
  createCluster: protectedProcedure
    .input(
      z.object({
        name: z.string(),
        organizationId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await clusterApiFactory
        .createCluster(userId, {
          name: input.name,
          organizationId: input.organizationId,
        })
        .then((res) => res.data);
    }),
  getCluster: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await clusterApiFactory
        .getCluster(userId, input.id)
        .then((res) => res.data);
    }),
  deleteCluster: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await clusterApiFactory
        .deleteCluster(userId, input.id)
        .then((res) => res.data);
    }),
};

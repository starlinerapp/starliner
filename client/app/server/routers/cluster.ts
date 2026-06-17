import { z } from "zod";
import { clusterApiFactory } from "~/server/api/clients/server";
import type { RequestCreateClusterServerTypeEnum } from "~/server/api/clients/server/generated";
import { streamSse } from "~/server/services/sse";
import { protectedProcedure } from "~/server/trpc";

export const clusterRouter = {
  createCluster: protectedProcedure
    .input(
      z.object({
        name: z.string(),
        serverType: z.string(),
        organizationId: z.number(),
        teamId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await clusterApiFactory
        .createCluster(userId, {
          name: input.name,
          serverType: input.serverType as RequestCreateClusterServerTypeEnum,
          organizationId: input.organizationId,
          teamId: input.teamId,
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
  streamProvisioningLogs: protectedProcedure
    .input(
      z.object({
        clusterId: z.number(),
      }),
    )
    .subscription(async function* ({ input, ctx, signal }) {
      const userId = ctx.user?.id;

      yield* streamSse(
        () =>
          // @ts-expect-error OpenAPI doesn't support SSE
          clusterApiFactory.streamClusterProvisioningLogs(
            userId,
            input.clusterId,
            { responseType: "stream", signal },
          ),
        signal,
        { trim: true },
      );
    }),
};

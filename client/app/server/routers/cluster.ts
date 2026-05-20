import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { clusterApiFactory } from "~/server/api/clients/server";
import type { RequestCreateClusterServerTypeEnum } from "~/server/api/clients/server/generated";
import type { AxiosResponse } from "axios";
import { Readable } from "stream";

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

      let response: AxiosResponse<Readable> | undefined;
      try {
        // @ts-expect-error OpenAPI doesn't support SSE
        response = await clusterApiFactory.streamClusterProvisioningLogs(
          userId,
          input.clusterId,
          { responseType: "stream", signal },
        );

        signal?.addEventListener("abort", () => {
          response?.data?.destroy();
        });

        const decoder = new TextDecoder();
        let buffer = "";

        // @ts-expect-error OpenAPI doesn't support SSE
        for await (const chunk of response.data) {
          if (signal?.aborted) {
            break;
          }

          buffer += decoder.decode(chunk, { stream: true });
          const lines = buffer.split("\n");
          buffer = lines.pop() ?? "";

          for (const line of lines) {
            if (line.startsWith("data: ")) {
              yield line.slice(6).trim();
            }
          }
        }
      } catch (err) {
        if (signal?.aborted) {
          return;
        }

        throw err;
      } finally {
        response?.data.destroy();
      }
    }),
};

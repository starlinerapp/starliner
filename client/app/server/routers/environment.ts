import { randomUUID } from "node:crypto";
import type { Readable } from "node:stream";
import type { AxiosResponse } from "axios";
import { z } from "zod";
import { environmentApiFactory } from "~/server/api/clients/server";
import { cache } from "~/server/services/cache";
import { protectedProcedure } from "~/server/trpc";

export const environmentRouter = {
  createEnvironment: protectedProcedure
    .input(
      z.object({
        name: z.string(),
        organizationId: z.number(),
        projectId: z.number(),
        sourceEnvironmentId: z.number().optional(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .createEnvironment(userId, {
          name: input.name,
          organizationId: input.organizationId,
          projectId: input.projectId,
          sourceEnvironmentId: input.sourceEnvironmentId,
        })
        .then((res) => res.data);
    }),
  deleteEnvironment: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .deleteEnvironment(userId, input.id)
        .then((res) => res.data);
    }),
  getEnvironmentDeployments: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .getEnvironmentDeployments(userId, input.id)
        .then((res) => res.data);
    }),
  getEnvironmentBuilds: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .getEnvironmentBuilds(userId, input.id)
        .then((res) => res.data);
    }),
  getEnvironmentConnectedBranch: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .getEnvironmentConnectedBranch(userId, input.id)
        .then((res) => res.data);
    }),
  updateEnvironmentConnectedBranch: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        branchName: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .updateEnvironmentConnectedBranch(userId, input.id, {
          branch: input.branchName,
        })
        .then((res) => res.data);
    }),
  streamDeploymentNotifications: protectedProcedure
    .input(
      z.object({
        id: z.number(),
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

      let response: AxiosResponse<Readable> | undefined;
      try {
        // @ts-expect-error OpenAPI doesn't support SSE
        response = await environmentApiFactory.streamEnvironmentNotifications(
          userId,
          correlationId,
          input.id,
          { responseType: "stream", signal },
        );

        signal?.addEventListener("abort", () => {
          response?.data?.destroy();
        });

        const decoder = new TextDecoder();
        let buffer = "";

        for await (const chunk of response!.data) {
          if (signal?.aborted) {
            break;
          }

          buffer += decoder.decode(chunk, { stream: true });
          const lines = buffer.split("\n");
          buffer = lines.pop() ?? "";

          for (const line of lines) {
            if (line.startsWith("data: ")) {
              const raw = line.slice(6).trim();
              if (raw) {
                try {
                  yield JSON.parse(raw);
                } catch {
                  yield raw;
                }
              }
            }
          }
        }
      } catch (e) {
        if (signal?.aborted) return;
        throw e;
      } finally {
        response?.data?.destroy();
      }
    }),
};

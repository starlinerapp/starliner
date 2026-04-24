import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { buildApiFactory } from "~/server/api/client";
import type { AxiosResponse } from "axios";
import { Readable } from "stream";

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

      let response: AxiosResponse<Readable> | undefined;
      try {
        // @ts-expect-error OpenAPI doesn't support SSE
        response = await buildApiFactory.streamBuildLogs(
          userId,
          input.buildId,
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

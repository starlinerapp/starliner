import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import type { AxiosResponse } from "axios";
import type { Readable } from "stream";
import { notificationsApiFactory } from "~/server/api/clients/server";

export const notificationsRouter = {
  streamGlobalNotifications: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .subscription(async function* ({ input, ctx, signal }) {
      const userId = ctx.user?.id;

      let response: AxiosResponse<Readable> | undefined;

      try {
        // @ts-expect-error OpenAPI doesn't support SSE
        response = await notificationsApiFactory.streamGlobalNotifications(
          userId,
          input.organizationId,
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

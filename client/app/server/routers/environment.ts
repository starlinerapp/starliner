import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { environmentApiFactory } from "~/api/client";
import { withAuthHeader } from "~/api/client/axios.server";

export const environmentRouter = {
  createEnvironment: protectedProcedure
    .input(
      z.object({
        name: z.string(),
        organizationId: z.number(),
        projectId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .createEnvironment(
          {
            name: input.name,
            organization_id: input.organizationId,
            project_id: input.projectId,
          },
          withAuthHeader(userId),
        )
        .then((res) => res.data);
    }),
};

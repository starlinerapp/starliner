import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { environmentApiFactory } from "~/server/api/client";

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
        .createEnvironment(userId, {
          name: input.name,
          organization_id: input.organizationId,
          project_id: input.projectId,
        })
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
};

import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { projectApiFactory } from "~/api/client";
import { withAuthHeader } from "~/api/client/axios.server";

export const projectRouter = {
  createProject: protectedProcedure
    .input(
      z.object({
        name: z.string(),
        organizationId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .createProject(
          {
            name: input.name,
            organization_id: input.organizationId,
          },
          withAuthHeader(userId),
        )
        .then((res) => res.data);
    }),
  getProject: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .getProject(input.id, withAuthHeader(userId))
        .then((res) => res.data);
    }),
};

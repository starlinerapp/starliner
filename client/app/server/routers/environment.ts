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
        sourceEnvironmentId: z.number().optional(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await environmentApiFactory
        .createEnvironment(userId, {
          name: input.name,
          organization_id: input.organizationId,
          project_id: input.projectId,
          source_environment_id: input.sourceEnvironmentId,
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
};

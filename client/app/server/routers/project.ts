import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { projectApiFactory } from "~/server/api/client";

export const projectRouter = {
  createProject: protectedProcedure
    .input(
      z.object({
        name: z.string(),
        clusterId: z.number(),
        teamId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .createProject(userId, {
          name: input.name,
          cluster_id: input.clusterId,
          team_id: input.teamId,
        })
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
        .getProject(userId, input.id)
        .then((res) => res.data);
    }),
  deleteProject: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .deleteProject(userId, input.id)
        .then((res) => res.data);
    }),
  getProjectCluster: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .getProjectCluster(userId, input.id)
        .then((res) => res.data);
    }),
  getProjectPreviewEnvironmentEnabled: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .getProjectPreviewEnvironmentEnabled(userId, input.id)
        .then((res) => res.data);
    }),
  toggleProjectPreviewEnvironmentEnabled: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await projectApiFactory
        .toggleProjectPreviewEnvironmentEnabled(userId, input.id)
        .then((res) => res.data);
    }),
};

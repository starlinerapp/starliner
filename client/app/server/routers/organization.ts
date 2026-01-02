import { protectedProcedure } from "~/server/trpc";
import { organizationApiFactory } from "~/server/api/client";
import { z } from "zod";

export const organizationRouter = {
  createOrganization: protectedProcedure
    .input(
      z.object({
        name: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .createOrganization(userId, { name: input.name })
        .then((res) => res.data);
    }),
  getUserOrganizations: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await organizationApiFactory
      .getUserOrganizations(userId)
      .then((res) => res.data);
  }),
  getOrganizationProjects: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .getOrganizationProjects(userId, input.id)
        .then((res) => res.data);
    }),
  getOrganizationClusters: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .getOrganizationClusters(userId, input.id)
        .then((res) => res.data);
    }),
};

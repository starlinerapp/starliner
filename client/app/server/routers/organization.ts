import { protectedProcedure } from "~/server/trpc";
import { organizationApiFactory } from "~/server/api/client";
import { withAuthHeader } from "~/server/api/client/axios.server";
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
        .createOrganization({ name: input.name }, withAuthHeader(userId))
        .then((res) => res.data);
    }),
  getUserOrganizations: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await organizationApiFactory
      .getUserOrganizations(withAuthHeader(userId))
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
        .getOrganizationProjects(input.id, withAuthHeader(userId))
        .then((res) => res.data);
    }),
};

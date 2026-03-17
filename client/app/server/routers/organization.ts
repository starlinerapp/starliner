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
  getUserProjects: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .getUserProjects(userId, input.id)
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
  upsertHetznerCredential: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        apiKey: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .upsertHetznerCredential(userId, input.id, {
          apiKey: input.apiKey,
        })
        .then((res) => res.data);
    }),
  getHetznerCredential: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .getHetznerCredential(userId, input.id)
        .then((res) => res.data);
    }),
  createInvite: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .createOrganizationInvite(userId, input.organizationId)
        .then((res) => res.data);
    }),
  acceptInvite: protectedProcedure
    .input(
      z.object({
        inviteId: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .acceptOrganizationInvite(userId, { inviteId: input.inviteId })
        .then((res) => res.data);
    }),
};

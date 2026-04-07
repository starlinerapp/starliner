import { protectedProcedure } from "~/server/trpc";
import { organizationApiFactory } from "~/server/api/client";
import { z } from "zod";
import { db } from "~/db";
import { user } from "~/db/schema";
import { inArray } from "drizzle-orm";

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
  getInvite: protectedProcedure
    .input(
      z.object({
        inviteId: z.string(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await organizationApiFactory
        .getOrganizationInviteDetails(userId, input.inviteId)
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
  getOrganizationMembers: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      const members = await organizationApiFactory
        .getOrganizationMembers(userId, input.id)
        .then((res) => res.data);

      const betterAuthIds = members.map((m) => m.better_auth_id);
      if (betterAuthIds.length === 0) return [];

      const authUsers = await db
        .select({ id: user.id, name: user.name, email: user.email })
        .from(user)
        .where(inArray(user.id, betterAuthIds));

      const authUserMap = new Map(authUsers.map((u) => [u.id, u]));

      return members.map((m) => ({
        ...m,
        name: authUserMap.get(m.better_auth_id)?.name ?? "",
        email: authUserMap.get(m.better_auth_id)?.email ?? "",
      }));
    }),
};

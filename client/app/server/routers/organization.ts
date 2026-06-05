import { TRPCError } from "@trpc/server";
import { isAxiosError } from "axios";
import { z } from "zod";
import { organizationApiFactory } from "~/server/api/clients/server";
import { enrichMembersWithAuthDetails } from "~/server/services/users";
import { protectedProcedure } from "~/server/trpc";

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
  sendInvite: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        toEmails: z.array(z.email()).min(1),
        inviteUrlPrefix: z.string().startsWith("/"),
        teamId: z.number().optional(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      const clientBaseUrl = process.env.CLIENT_BASE_URL;

      if (!clientBaseUrl) {
        throw new Error("Environment variable 'CLIENT_BASE_URL' is not set");
      }

      try {
        return await organizationApiFactory
          .sendOrganizationInvite(userId, input.organizationId, {
            toEmails: input.toEmails,
            inviteUrlPrefix: clientBaseUrl + input.inviteUrlPrefix,
            teamId: input.teamId,
          })
          .then((res) => res.data);
      } catch (err) {
        if (isAxiosError(err) && err.response?.data?.error) {
          throw new TRPCError({
            code:
              err.response.status === 400
                ? "BAD_REQUEST"
                : "INTERNAL_SERVER_ERROR",
            message: String(err.response.data.error),
          });
        }
        throw err;
      }
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
        .acceptOrganizationInvite(userId, {
          inviteId: input.inviteId,
          recipientEmail: ctx.user?.email,
        })
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

      return await enrichMembersWithAuthDetails(members);
    }),
  removeOrganizationMember: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        userId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const callerId = ctx.user?.id;
      return await organizationApiFactory
        .removeOrganizationMember(callerId, input.organizationId, {
          userId: input.userId,
        })
        .then((res) => res.data);
    }),
};

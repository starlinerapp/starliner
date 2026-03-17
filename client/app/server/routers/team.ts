import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { teamsApiFactory } from "~/server/api/client";

export const teamRouter = {
  createTeam: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        name: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .createTeam(userId, input.organizationId, {
          name: input.name,
        })
        .then((res) => res.data);
    }),
  getUserTeams: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .getUserTeams(userId, input.organizationId)
        .then((res) => res.data);
    }),
  getTeamMembers: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        teamId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .getTeamMembers(userId, input.organizationId, input.teamId)
        .then((res) => res.data);
    }),
  joinTeam: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        slug: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .joinTeam(userId, input.organizationId, {
          slug: input.slug,
        })
        .then((res) => res.data);
    }),
  removeTeamMember: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        teamId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .removeTeamMember(userId, input.organizationId, input.teamId)
        .then((res) => res.data);
    }),
};

import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { teamsApiFactory } from "~/server/api/client";
import { db } from "~/db";
import { user } from "~/db/schema";
import { inArray } from "drizzle-orm";

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
          slug: input.name,
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
        teamId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      const members = await teamsApiFactory
        .getTeamMembers(userId, input.teamId)
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
  addTeamMember: protectedProcedure
    .input(
      z.object({
        teamId: z.number(),
        userId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const callerId = ctx.user?.id;
      return await teamsApiFactory
        .addTeamMember(callerId, input.teamId, { userId: input.userId })
        .then((res) => res.data);
    }),
  removeTeamMember: protectedProcedure
    .input(
      z.object({
        teamId: z.number(),
        userId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const callerId = ctx.user?.id;
      return await teamsApiFactory
        .removeTeamMember(callerId, input.teamId, { userId: input.userId })
        .then((res) => res.data);
    }),
  getTeamRepositories: protectedProcedure
    .input(
      z.object({
        teamId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .getTeamRepositories(userId, input.teamId)
        .then((res) => res.data);
    }),
  assignRepoToTeam: protectedProcedure
    .input(
      z.object({
        teamId: z.number(),
        githubRepoId: z.number(),
        repoName: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .assignRepoToTeam(userId, input.teamId, {
          githubRepoId: input.githubRepoId,
          repoName: input.repoName,
        })
        .then((res) => res.data);
    }),
  unassignRepoFromTeam: protectedProcedure
    .input(
      z.object({
        teamId: z.number(),
        githubRepoId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await teamsApiFactory
        .unassignRepoFromTeam(userId, input.teamId, input.githubRepoId)
        .then((res) => res.data);
    }),
};

import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { githubFactory } from "~/server/api/client";

export const githubRouter = {
  getRepositories: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      try {
        const res = await githubFactory.getRepositories(
          userId,
          input.organizationId,
        );

        return res.data;
      } catch (err) {
        console.error(err);
        throw err;
      }
    }),
  getAllRepositories: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      try {
        const res = await githubFactory.getAllRepositories(
          userId,
          input.organizationId,
        );

        return res.data;
      } catch (err) {
        console.error(err);
        throw err;
      }
    }),
  getRepositoryFiles: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        owner: z.string(),
        repo: z.string(),
        path: z.string().optional(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      const res = await githubFactory.getRepositoryContents(
        userId,
        input.organizationId,
        input.owner,
        input.repo,
        input.path,
      );
      return res.data;
    }),
};

import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { githubApiFactory } from "~/server/api/clients/server";

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
        const res = await githubApiFactory.getRepositories(
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
        const res = await githubApiFactory.getAllRepositories(
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
      const res = await githubApiFactory.getRepositoryContents(
        userId,
        input.organizationId,
        input.owner,
        input.repo,
        input.path,
      );
      return res.data;
    }),
  getRepositoryFileContent: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
        owner: z.string(),
        repo: z.string(),
        path: z.string(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await githubApiFactory
        .getFileContent(
          userId,
          input.organizationId,
          input.owner,
          input.repo,
          input.path,
        )
        .then((res) => res.data);
    }),
};

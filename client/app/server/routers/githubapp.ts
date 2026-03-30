import { protectedProcedure } from "~/server/trpc";
import { githubAppFactory } from "~/server/api/client";
import z from "zod";
import axios from "axios";

export const githubAppRouter = {
  createGithubApp: protectedProcedure
    .input(
      z.object({
        installationId: z.number(),
        organizationId: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await githubAppFactory
        .createGithubApp(userId, {
          installationId: input.installationId,
          organizationId: input.organizationId,
        })
        .then((res) => res.data);
    }),
  getGithubApp: protectedProcedure
    .input(
      z.object({
        organizationId: z.number(),
      }),
    )
    .query(async ({ input, ctx }) => {
      const userId = ctx.user?.id;

      try {
        const res = await githubAppFactory.getGithubApp(
          userId,
          input.organizationId,
        );
        return res.data;
      } catch (err: unknown) {
        if (axios.isAxiosError(err) && err.response?.status === 404) {
          return null;
        }
        throw err;
      }
    }),
};

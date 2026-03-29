import { protectedProcedure } from "~/server/trpc";
import { githubAppFactory } from "~/server/api/client";
import z from "zod";

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
};

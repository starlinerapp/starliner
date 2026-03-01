import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { deploymentApiFactory } from "~/server/api/client";

const ingressPathSchema = z.object({
  path: z.string(),
  pathType: z.enum(["Prefix", "Exact"]),
  serviceName: z.string(),
});

const ingressHostSchema = z.object({
  host: z.string(),
  paths: z.array(ingressPathSchema),
});

export const deploymentRouter = {
  deployImage: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        serviceName: z.string(),
        imageName: z.string(),
        tag: z.string(),
        port: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .deployImage(userId, {
          environmentId: input.id,
          serviceName: input.serviceName,
          imageName: input.imageName,
          tag: input.tag,
          port: input.port,
        })
        .then((res) => res.data);
    }),
  deployDatabase: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .deployDatabase(userId, {
          environmentId: input.id,
          database: "postgres",
        })
        .then((res) => res.data);
    }),
  deleteDeployment: protectedProcedure
    .input(
      z.object({
        id: z.number(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .deleteDeployment(userId, input.id)
        .then((res) => res.data);
    }),
  deployIngress: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        ingressHosts: z.array(ingressHostSchema),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .deployIngress(userId, {
          environmentId: input.id,
          ingressHosts: input.ingressHosts,
        })
        .then((res) => res.data);
    }),
};

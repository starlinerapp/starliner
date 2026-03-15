import { protectedProcedure } from "~/server/trpc";
import { z } from "zod";
import { deploymentApiFactory } from "~/server/api/client";
import { type AxiosResponse, isAxiosError } from "axios";
import { Readable } from "stream";
import { TRPCError } from "@trpc/server";

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
  deployFromGitRepo: protectedProcedure
    .input(
      z.object({
        environmentId: z.number(),
        serviceName: z.string(),
        port: z.number(),
        gitUrl: z.string(),
        dockerfilePath: z.string(),
        projectRepositoryPath: z.string(),
        envs: z
          .array(
            z.object({
              name: z.string(),
              value: z.string(),
            }),
          )
          .default([]),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      try {
        const res = await deploymentApiFactory.deployFromGitRepository(userId, {
          environmentId: input.environmentId,
          serviceName: input.serviceName,
          port: input.port,
          gitUrl: input.gitUrl,
          dockerfilePath: input.dockerfilePath,
          projectRepositoryPath: input.projectRepositoryPath,
          envs: input.envs,
        });
        return res.data;
      } catch (err) {
        if (isAxiosError(err) && err.response?.data?.error) {
          throw new TRPCError({
            code:
              err.response.status === 409
                ? "CONFLICT"
                : "INTERNAL_SERVER_ERROR",
            message: err.response.data.error,
          });
        }
        throw err;
      }
    }),
  updateDeployFromGitRepo: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        deploymentId: z.number(),
        port: z.number(),
        dockerfilePath: z.string(),
        projectRepositoryPath: z.string(),
        envs: z
          .array(
            z.object({
              name: z.string(),
              value: z.string(),
            }),
          )
          .default([]),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .updateDeployFromGitRepository(userId, input.deploymentId, {
          environmentId: input.id,
          port: input.port,
          dockerfilePath: input.dockerfilePath,
          projectRepositoryPath: input.projectRepositoryPath,
          envs: input.envs,
        })
        .then((res) => res.data);
    }),
  deployImage: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        serviceName: z.string(),
        imageName: z.string(),
        tag: z.string(),
        port: z.number(),
        envs: z
          .array(
            z.object({
              name: z.string(),
              value: z.string(),
            }),
          )
          .default([]),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      try {
        const res = await deploymentApiFactory.deployImage(userId, {
          environmentId: input.id,
          serviceName: input.serviceName,
          imageName: input.imageName,
          tag: input.tag,
          port: input.port,
          envs: input.envs,
        });
        return res.data;
      } catch (err) {
        if (isAxiosError(err) && err.response?.data?.error) {
          throw new TRPCError({
            code:
              err.response.status === 409
                ? "CONFLICT"
                : "INTERNAL_SERVER_ERROR",
            message: err.response.data.error,
          });
        }
        throw err;
      }
    }),
  updateImage: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        deploymentId: z.number(),
        imageName: z.string(),
        tag: z.string(),
        port: z.number(),
        envs: z
          .array(
            z.object({
              name: z.string(),
              value: z.string(),
            }),
          )
          .default([]),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .updateImageDeployment(userId, input.deploymentId, {
          environmentId: input.id,
          imageName: input.imageName,
          tag: input.tag,
          port: input.port,
          envs: input.envs,
        })
        .then((res) => res.data);
    }),
  deployDatabase: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        serviceName: z.string(),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      try {
        const res = await deploymentApiFactory.deployDatabase(userId, {
          environmentId: input.id,
          serviceName: input.serviceName,
        });
        return res.data;
      } catch (err) {
        if (isAxiosError(err) && err.response?.data?.error) {
          throw new TRPCError({
            code:
              err.response.status === 409
                ? "CONFLICT"
                : "INTERNAL_SERVER_ERROR",
            message: err.response.data.error,
          });
        }
        throw err;
      }
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
  updateIngress: protectedProcedure
    .input(
      z.object({
        id: z.number(),
        deploymentId: z.number(),
        ingressHosts: z.array(ingressHostSchema),
      }),
    )
    .mutation(async ({ input, ctx }) => {
      const userId = ctx.user?.id;
      return await deploymentApiFactory
        .updateIngressDeployment(userId, input.deploymentId, {
          environmentId: input.id,
          ingressHosts: input.ingressHosts,
        })
        .then((res) => res.data);
    }),
  streamDeploymentLogs: protectedProcedure
    .input(
      z.object({
        deploymentId: z.number(),
      }),
    )
    .subscription(async function* ({ input, ctx, signal }) {
      const userId = ctx.user?.id;

      let response: AxiosResponse<Readable> | undefined;
      try {
        // @ts-expect-error OpenAPI doesn't support SSE
        response = await deploymentApiFactory.streamDeploymentLogs(
          userId,
          input.deploymentId,
          { responseType: "stream", signal },
        );

        signal?.addEventListener("abort", () => {
          response?.data?.destroy();
        });

        const decoder = new TextDecoder();
        let buffer = "";

        // @ts-expect-error OpenAPI doesn't support SSE
        for await (const chunk of response.data) {
          if (signal?.aborted) {
            break;
          }

          buffer += decoder.decode(chunk, { stream: true });
          const lines = buffer.split("\n");
          buffer = lines.pop() ?? "";

          for (const line of lines) {
            if (line.startsWith("data: ")) {
              yield line.slice(6).trim();
            }
          }
        }
      } catch (err) {
        if (signal?.aborted) {
          return;
        }

        throw err;
      } finally {
        response?.data.destroy();
      }
    }),
};

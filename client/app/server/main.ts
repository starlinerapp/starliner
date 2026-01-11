import { createTRPCRouter } from "~/server/trpc";
import { rootRouter } from "~/server/routers/root";
import { userRouter } from "~/server/routers/user";
import { organizationRouter } from "~/server/routers/organization";
import { projectRouter } from "~/server/routers/project";
import { environmentRouter } from "~/server/routers/environment";
import { clusterRouter } from "~/server/routers/cluster";
import { deploymentRouter } from "~/server/routers/deployment";

export const appRouter = createTRPCRouter({
  root: rootRouter,
  user: userRouter,
  organization: organizationRouter,
  project: projectRouter,
  environment: environmentRouter,
  cluster: clusterRouter,
  deployment: deploymentRouter,
});

export type AppRouter = typeof appRouter;

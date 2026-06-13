import { buildRouter } from "~/server/routers/build";
import { clusterRouter } from "~/server/routers/cluster";
import { deploymentRouter } from "~/server/routers/deployment";
import { environmentRouter } from "~/server/routers/environment";
import { githubRouter } from "~/server/routers/github";
import { githubAppRouter } from "~/server/routers/githubapp";
import { notificationsRouter } from "~/server/routers/notifications";
import { organizationRouter } from "~/server/routers/organization";
import { projectRouter } from "~/server/routers/project";
import { rootRouter } from "~/server/routers/root";
import { teamRouter } from "~/server/routers/team";
import { userRouter } from "~/server/routers/user";
import { createTRPCRouter } from "~/server/trpc";

export const appRouter = createTRPCRouter({
  root: rootRouter,
  user: userRouter,
  organization: organizationRouter,
  project: projectRouter,
  environment: environmentRouter,
  cluster: clusterRouter,
  deployment: deploymentRouter,
  build: buildRouter,
  team: teamRouter,
  github: githubRouter,
  githubApp: githubAppRouter,
  notifications: notificationsRouter,
});

export type AppRouter = typeof appRouter;

import { createTRPCRouter } from "~/server/trpc";
import { rootRouter } from "~/server/routers/root";
import { userRouter } from "~/server/routers/user";
import { organizationRouter } from "~/server/routers/organization";
import { projectRouter } from "~/server/routers/project";
import { environmentRouter } from "~/server/routers/environment";
import { clusterRouter } from "~/server/routers/cluster";
import { deploymentRouter } from "~/server/routers/deployment";
import { buildRouter } from "~/server/routers/build";
import { teamRouter } from "~/server/routers/team";
import { githubAppRouter } from "~/server/routers/githubapp";
import { githubRouter } from "~/server/routers/github";

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
});

export type AppRouter = typeof appRouter;

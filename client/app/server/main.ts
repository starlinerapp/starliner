import { createTRPCRouter } from "~/server/trpc";
import { rootRouter } from "~/server/routers/root";
import { userRouter } from "~/server/routers/user";
import { organizationRouter } from "~/server/routers/organization";
import { projectRouter } from "~/server/routers/project";
import { environmentRouter } from "~/server/routers/environment";

export const appRouter = createTRPCRouter({
  root: rootRouter,
  user: userRouter,
  organization: organizationRouter,
  project: projectRouter,
  environment: environmentRouter,
});

export type AppRouter = typeof appRouter;

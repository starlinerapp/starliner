import { createTRPCRouter } from "~/server/trpc";
import { rootRouter } from "~/server/routers/root";

export const appRouter = createTRPCRouter({
  root: rootRouter,
});

export type AppRouter = typeof appRouter;

import { createTRPCRouter } from "~/server/trpc";
import { rootRouter } from "~/server/routers/root";
import { userRouter } from "~/server/routers/user";

export const appRouter = createTRPCRouter({
  root: rootRouter,
  user: userRouter,
});

export type AppRouter = typeof appRouter;

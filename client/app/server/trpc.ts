import superjson from "superjson";
import { initTRPC, TRPCError } from "@trpc/server";
import { ZodError } from "zod";
import { auth } from "~/utils/auth/server";

export async function createTRPCContext(opts: { headers: Headers }) {
  const authSession = await auth.api.getSession({
    headers: opts.headers,
  });
  return {
    session: authSession?.session,
    user: authSession?.user,
  };
}

export type Context = Awaited<ReturnType<typeof createTRPCContext>>;

const trpc = initTRPC.context<Context>().create({
  transformer: superjson,
  errorFormatter: ({ shape, error }) => ({
    ...shape,
    data: {
      ...shape.data,
      zodError: error.cause instanceof ZodError ? error.cause : null,
    },
  }),
});

// Caller Factory for making server-side tRPC calls from loaders or actions
export const createCallerFactory = trpc.createCallerFactory;
export const createTRPCRouter = trpc.router;
export const publicProcedure = trpc.procedure;

export const protectedProcedure = trpc.procedure.use(({ ctx, next }) => {
  if (!ctx.user?.id) {
    throw new TRPCError({ code: "UNAUTHORIZED" });
  }
  return next({
    ctx: {
      user: ctx.user,
    },
  });
});

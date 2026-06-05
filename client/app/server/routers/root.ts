import { rootApiFactory } from "~/server/api/clients/server";
import { protectedProcedure } from "~/server/trpc";

export const rootRouter = {
  getRoot: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await rootApiFactory.getRoot(userId).then((res) => res.data);
  }),
};

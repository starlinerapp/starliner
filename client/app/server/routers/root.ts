import { protectedProcedure } from "~/server/trpc";
import { rootApiFactory } from "~/api/client";
import { withAuthHeader } from "~/api/client/axios.server";

export const rootRouter = {
  getRoot: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await rootApiFactory
      .getRoot(withAuthHeader(userId))
      .then((res) => res.data);
  }),
};

import { protectedProcedure } from "~/server/trpc";
import { userApiFactory } from "~/api/client";
import { withAuthHeader } from "~/api/client/axios.server";

export const userRouter = {
  getUser: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await userApiFactory
      .getUser(withAuthHeader(userId))
      .then((res) => res.data);
  }),
};

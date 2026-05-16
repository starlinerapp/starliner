import { protectedProcedure } from "~/server/trpc";
import { userApiFactory } from "~/server/api/clients/server";

export const userRouter = {
  getUser: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await userApiFactory.getUser(userId).then((res) => res.data);
  }),
};

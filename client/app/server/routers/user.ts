import { userApiFactory } from "~/server/api/clients/server";
import { protectedProcedure } from "~/server/trpc";

export const userRouter = {
  getUser: protectedProcedure.query(async ({ ctx }) => {
    const userId = ctx.user?.id;
    return await userApiFactory.getUser(userId).then((res) => res.data);
  }),
};

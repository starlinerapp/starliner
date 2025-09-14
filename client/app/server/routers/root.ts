import { publicProcedure } from "~/server/trpc";
import { rootApiFactory } from "~/api/client";

export const rootRouter = {
  getRoot: publicProcedure.query(async () => {
    return await rootApiFactory.getRoot().then((res) => res.data);
  }),
};

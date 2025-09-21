import { Configuration, RootApiFactory } from "~/api/client/generated";
import { serverEnv } from "~/env.server";
import { axiosInstance } from "~/api/client/axios.server";

const configuration = new Configuration({
  basePath: `http://${serverEnv.SERVER_BASE_URL}`,
});

export const rootApiFactory = RootApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

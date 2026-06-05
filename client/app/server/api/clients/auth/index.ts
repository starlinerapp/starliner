import axios from "axios";
import { serverEnv } from "~/env.server";
import {
  Configuration,
  UsersApiFactory,
} from "~/server/api/clients/auth/generated";

const axiosInstance = axios.create();

const configuration = new Configuration({
  basePath: serverEnv.AUTH_URL,
});

export const usersApiFactory = UsersApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

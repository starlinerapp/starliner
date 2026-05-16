import {
  Configuration,
  UsersApiFactory,
} from "~/server/api/clients/auth/generated";
import { serverEnv } from "~/env.server";
import axios from "axios";

const axiosInstance = axios.create();

const configuration = new Configuration({
  basePath: serverEnv.AUTH_URL,
});

export const usersApiFactory = UsersApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

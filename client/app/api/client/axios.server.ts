import axios, { type RawAxiosRequestConfig } from "axios";
import type { Id } from "~/types";
import { serverEnv } from "~/env.server";

export const axiosInstance = axios.create();

axiosInstance.interceptors.request.use((config) => {
  if (serverEnv.SERVER_BASE_URL) {
    config.auth = {
      username: serverEnv.SERVER_BASIC_AUTH_USER,
      password: serverEnv.SERVER_BASIC_AUTH_PASSWORD,
    };
  }

  return config;
});

export const withAuthHeader = (userId: Id, options?: RawAxiosRequestConfig) => {
  const reqConfig = options ?? {};

  reqConfig.headers = reqConfig.headers ?? {};
  reqConfig.headers["X-User-ID"] = userId;
  return reqConfig;
};

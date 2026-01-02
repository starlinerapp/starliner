import axios from "axios";
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

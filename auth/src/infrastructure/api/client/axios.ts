import axios from "axios";
import { serverEnv } from "../../../env.server";

export function createServerApiAxios() {
  const instance = axios.create();

  instance.interceptors.request.use((config) => {
    config.auth = {
      username: serverEnv.SERVER_BASIC_AUTH_USER,
      password: serverEnv.SERVER_BASIC_AUTH_PASSWORD,
    };

    return config;
  });

  return instance;
}

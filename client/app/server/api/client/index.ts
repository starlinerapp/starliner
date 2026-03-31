import {
  OrganizationApiFactory,
  Configuration,
  RootApiFactory,
  UserApiFactory,
  ProjectApiFactory,
  EnvironmentApiFactory,
  ClusterApiFactory,
  DeploymentApiFactory,
  BuildApiFactory,
  TeamApiFactory,
  GithubappApiFactory,
  GithubApiFactory,
} from "~/server/api/client/generated";

import { serverEnv } from "~/env.server";
import { axiosInstance } from "~/server/api/client/axios.server";

const configuration = new Configuration({
  basePath: `http://${serverEnv.SERVER_BASE_URL}`,
});

export const rootApiFactory = RootApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const userApiFactory = UserApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const organizationApiFactory = OrganizationApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const projectApiFactory = ProjectApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const environmentApiFactory = EnvironmentApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const clusterApiFactory = ClusterApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const deploymentApiFactory = DeploymentApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const buildApiFactory = BuildApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const teamsApiFactory = TeamApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const githubFactory = GithubApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

export const githubAppFactory = GithubappApiFactory(
  configuration,
  undefined,
  axiosInstance,
);

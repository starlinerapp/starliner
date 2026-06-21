export type EnvironmentBuildPollState = {
  status: string;
  deploymentRolloutStatus: string;
};

export function shouldPollEnvironmentBuilds(
  builds: EnvironmentBuildPollState[] | undefined,
): number | false {
  if (!builds) {
    return false;
  }

  const shouldPoll = builds.some((build) => isEnvironmentBuildInProgress(build));

  return shouldPoll ? 1000 : false;
}

export function isEnvironmentBuildInProgress(
  build: EnvironmentBuildPollState,
): boolean {
  if (
    build.status === "failure" ||
    build.deploymentRolloutStatus === "failure"
  ) {
    return false;
  }

  if (
    build.status === "success" &&
    build.deploymentRolloutStatus === "success"
  ) {
    return false;
  }

  return (
    build.status === "building" ||
    build.status === "queued" ||
    (build.status === "success" && build.deploymentRolloutStatus === "pending")
  );
}

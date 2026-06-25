import DeploymentCard from "~/components/organisms/deployment-card/DeploymentCard";
import { useEnvironmentLayoutContext } from "../layout";

export default function Builds() {
  const { environmentBuilds, isEnvironmentBuildsLoading } =
    useEnvironmentLayoutContext();

  if (!isEnvironmentBuildsLoading && environmentBuilds?.length === 0)
    return (
      <div className="flex flex-col gap-1 p-4">
        <p className="text-mauve-11">
          There are no deployments for this environment yet.
        </p>
      </div>
    );

  return (
    <div className="flex flex-col gap-4 p-4">
      {environmentBuilds?.map((build, i) => (
        <DeploymentCard
          isCollapsed={i > 0}
          key={build.buildId}
          buildId={build.buildId}
          deploymentId={build.deploymentId}
          commitHash={build.commitHash}
          source={build.source}
          serviceName={build.deploymentName}
          createdAt={build.createdAt}
          status={build.status}
          deploymentRolloutStatus={build.deploymentRolloutStatus}
        />
      ))}
    </div>
  );
}

import React from "react";
import DeployDatabaseForm from "~/components/organisms/forms/DeployDatabaseForm";
import { useParams } from "react-router";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";

export default function UpdateDatabaseDeployment() {
  const { deploymentId } = useParams<{ deploymentId: string }>();

  const { environment: currentEnvironment } = useEnvironment();

  const trpc = useTRPC();
  const { data: environmentDeployments, isLoading } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: currentEnvironment?.id },
      { enabled: !!currentEnvironment },
    ),
  );

  const databaseDeployment = environmentDeployments?.databases.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployDatabaseForm
          key={deploymentId}
          defaultValues={{
            serviceName: databaseDeployment?.serviceName ?? "",
          }}
        />
      )}
    </>
  );
}

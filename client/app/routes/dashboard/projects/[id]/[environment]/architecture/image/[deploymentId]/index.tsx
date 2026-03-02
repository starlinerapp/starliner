import React from "react";
import DeployImageForm from "~/components/organisms/forms/DeployImageForm";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useParams } from "react-router";

export default function UpdateImageForm() {
  const { deploymentId } = useParams<{ deploymentId: string }>();

  const { environment: currentEnvironment } = useEnvironment();

  const trpc = useTRPC();
  const { data: environmentDeployments, isLoading } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions(
      { id: currentEnvironment?.id },
      { enabled: !!currentEnvironment },
    ),
  );

  const imageDeployment = environmentDeployments?.images.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployImageForm
          defaultValues={{
            serviceName: imageDeployment?.serviceName ?? "",
            imageName: imageDeployment?.imageName ?? "",
            tag: imageDeployment?.tag ?? "",
            port: imageDeployment ? Number(imageDeployment.port) : undefined,
            envs: imageDeployment?.envVars ?? [],
          }}
        />
      )}
    </>
  );
}

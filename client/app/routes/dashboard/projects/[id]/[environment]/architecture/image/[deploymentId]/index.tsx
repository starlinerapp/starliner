import React from "react";
import DeployImageForm, {
  type ImageFormInput,
} from "~/components/organisms/forms/DeployImageForm";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery } from "@tanstack/react-query";
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

  const updateImageMutation = useMutation(
    trpc.deployment.updateImage.mutationOptions(),
  );

  const onSubmit = async (data: ImageFormInput) => {
    if (!data.port) {
      return;
    }

    await updateImageMutation.mutateAsync({
      id: currentEnvironment.id,
      deploymentId: Number(deploymentId),
      imageName: data.imageName,
      tag: data.tag,
      port: data.port,
      envs: data.envs,
    });
  };

  const imageDeployment = environmentDeployments?.images.find(
    (deployment) => deployment.id === Number(deploymentId),
  );

  return (
    <>
      {!isLoading && (
        <DeployImageForm
          key={deploymentId}
          onSubmit={onSubmit}
          defaultValues={{
            serviceName: imageDeployment?.serviceName ?? "",
            imageName: imageDeployment?.imageName ?? "",
            tag: imageDeployment?.tag ?? "",
            port: imageDeployment ? Number(imageDeployment.port) : null,
            volumeSizeMiB: imageDeployment?.volumeSizeMiB ?? null,
            volumeMountPath: imageDeployment?.volumeMountPath ?? null,
            envs: imageDeployment?.envVars ?? [],
          }}
        />
      )}
    </>
  );
}

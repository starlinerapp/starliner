import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import DeployImageForm, {
  type ImageFormInput,
} from "~/components/organisms/forms/DeployImageForm";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";

export default function UpdateImageForm() {
  const { slug, id, environment, deploymentId } = useParams<{
    slug: string;
    id: string;
    environment: string;
    deploymentId: string;
  }>();
  const navigate = useNavigate();

  const { environment: currentEnvironment } = useEnvironment();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
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

    await updateImageMutation.mutateAsync(
      {
        id: currentEnvironment.id,
        deploymentId: Number(deploymentId),
        imageName: data.imageName,
        tag: data.tag,
        port: data.port,
        envs: data.envs,
      },
      {
        onSuccess: (result) => {
          queryClient.invalidateQueries({
            queryKey: trpc.environment.getEnvironmentBuilds.queryKey({
              id: currentEnvironment.id,
            }),
          });
          queryClient.invalidateQueries({
            queryKey: trpc.environment.getEnvironmentDeployments.queryKey({
              id: currentEnvironment.id,
            }),
          });
          if (result?.deploymentId && slug && id && environment) {
            navigate(
              `/${slug}/projects/${id}/${environment}/architecture/image/${result.deploymentId}`,
              { replace: true },
            );
          }
        },
      },
    );
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

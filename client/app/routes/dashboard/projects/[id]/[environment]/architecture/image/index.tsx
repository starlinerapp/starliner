import { useMutation, useQueryClient } from "@tanstack/react-query";
import DeployImageForm, {
  type ImageFormInput,
} from "~/components/organisms/forms/DeployImageForm";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useTRPC } from "~/utils/trpc/react";

export default function Index() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();

  const { environment: currentEnvironment } = useEnvironment();

  const createImageMutation = useMutation(
    trpc.deployment.deployImage.mutationOptions(),
  );

  const onSubmit = async (data: ImageFormInput) => {
    if (!data.port) {
      return;
    }

    await createImageMutation.mutateAsync(
      {
        id: currentEnvironment.id,
        serviceName: data.serviceName,
        imageName: data.imageName,
        tag: data.tag,
        port: data.port,
        volumeSizeMiB: data.volumeSizeMiB ?? undefined,
        volumeMountPath: data.volumeMountPath ?? undefined,
        envs: data.envs,
      },
      {
        onSuccess: () => {
          queryClient.invalidateQueries({
            queryKey: trpc.environment.getEnvironmentBuilds.queryKey({
              id: currentEnvironment.id,
            }),
          });
        },
      },
    );
  };

  return <DeployImageForm resetOnSuccess={true} onSubmit={onSubmit} />;
}

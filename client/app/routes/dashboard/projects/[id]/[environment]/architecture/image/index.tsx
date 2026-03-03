import React from "react";
import DeployImageForm, {
  type ImageFormInput,
} from "~/components/organisms/forms/DeployImageForm";
import { useMutation } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";

export default function Index() {
  const trpc = useTRPC();

  const { environment: currentEnvironment } = useEnvironment();

  const createImageMutation = useMutation(
    trpc.deployment.deployImage.mutationOptions(),
  );

  const onSubmit = async (data: ImageFormInput) => {
    if (!data.port) {
      return;
    }

    await createImageMutation.mutateAsync({
      id: currentEnvironment.id,
      serviceName: data.serviceName,
      imageName: data.imageName,
      tag: data.tag,
      port: data.port,
      envs: data.envs,
    });
  };

  return <DeployImageForm resetOnSuccess={true} onSubmit={onSubmit} />;
}

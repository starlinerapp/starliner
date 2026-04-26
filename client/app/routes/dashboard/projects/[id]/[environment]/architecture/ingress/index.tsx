import React from "react";
import DeployIngressForm, {
  type IngressFormInput,
} from "~/components/organisms/forms/DeployIngressForm";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";

export default function Index() {
  const trpc = useTRPC();
  const createIngressMutation = useMutation(
    trpc.deployment.deployIngress.mutationOptions(),
  );
  const { environment: currentEnvironment } = useEnvironment();

  const onSubmit = async (data: IngressFormInput) => {
    await createIngressMutation.mutateAsync({
      id: currentEnvironment.id,
      ingressHosts: data.hosts.map((h) => ({
        host: h.name,
        paths: h.paths.map((p) => ({
          path: p.path,
          pathType: p.pathType as "Prefix" | "Exact",
          serviceName: p.service,
        })),
      })),
    });
  };

  return <DeployIngressForm resetOnSuccess={true} onSubmit={onSubmit} />;
}

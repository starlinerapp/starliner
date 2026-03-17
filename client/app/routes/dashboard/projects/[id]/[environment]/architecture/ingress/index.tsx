import React from "react";
import DeployIngressForm, {
  type IngressFormInput,
} from "~/components/organisms/forms/DeployIngressForm";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function Index() {
  const trpc = useTRPC();
  const createIngressMutation = useMutation(
    trpc.deployment.deployIngress.mutationOptions(),
  );
  const { environment: currentEnvironment, clusterId } = useEnvironment();
  const { data: clusterData } = useQuery(
    trpc.cluster.getCluster.queryOptions(
      { id: clusterId! },
      { enabled: !!clusterId },
    ),
  );

  const organization = useOrganizationContext();

  const onSubmit = async (data: IngressFormInput) => {
    await createIngressMutation.mutateAsync({
      id: currentEnvironment.id,
      ingressHosts: data.hosts.map((h) => ({
        host:
          h.name + `.${organization.slug}.${clusterData?.ipv4Address}.nip.io`,
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

import React from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function Settings() {
  const trpc = useTRPC();
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const queryClient = useQueryClient();
  const { id } = useParams<{
    id: string;
  }>();

  const deleteClusterMutation = useMutation(
    trpc.cluster.deleteCluster.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationClusters.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${organization.slug}/clusters/all`);
      },
    }),
  );

  return (
    <div className="p-4">
      <div className="flex items-center gap-6 py-2">
        <div className="flex flex-col">
          <p className="text-md font-bold">Delete this Cluster</p>
          <p className="text-sm">
            Once you delete a cluster, there is no going back. Please be
            certain.
          </p>
        </div>
        <div>
          <Button
            intent="danger"
            size="sm"
            onClick={() =>
              deleteClusterMutation.mutate({
                id: Number(id),
              })
            }
          >
            Delete this Cluster
          </Button>
        </div>
      </div>
    </div>
  );
}

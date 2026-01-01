import React from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function Settings() {
  const trpc = useTRPC();
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
      },
    }),
  );

  return (
    <div className="w-3/5 p-4">
      <div className="border-mauve-6 rounded-md border-1">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-sm uppercase">
          Danger Zone
        </div>
        <div className="flex items-center justify-between p-4">
          <div>
            <p className="text-md font-bold">Delete this Cluster</p>
            <p className="text-mauve-11 text-sm">
              Once you delete a cluster, there is no going back. Please be
              certain.
            </p>
          </div>
          <Button
            className="w-36"
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

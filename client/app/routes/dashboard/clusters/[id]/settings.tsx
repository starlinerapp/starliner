import React, { useEffect } from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function Settings() {
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const { slug, id } = useParams<{
    slug: string;
    id: string;
  }>();

  const { data: clusterData, error } = useQuery(
    trpc.cluster.getCluster.queryOptions(
      { id: Number(id) },
      {
        refetchInterval: 4000,
      },
    ),
  );

  useEffect(() => {
    if (error) {
      (async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationClusters.queryKey({
            id: organization.id,
          }),
        });
        navigate(`clusters/all`);
      })();
    }
  }, [error]);

  const deleteClusterMutation = useMutation(
    trpc.cluster.deleteCluster.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.cluster.getCluster.queryKey({
            id: Number(id),
          }),
        });
        navigate(`/${slug}/clusters/${id}/general`);
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
            disabled={
              clusterData?.status === "pending" ||
              clusterData?.status === "deleted"
            }
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

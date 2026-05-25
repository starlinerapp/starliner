import React, { useEffect, useState } from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import DestructiveDialog from "~/components/organisms/dialog/DestructiveDialog";

export default function Settings() {
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const { slug, id } = useParams<{
    slug: string;
    id: string;
  }>();

  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const { data: hetznerCredentialData, isLoading: isCredentialLoading } =
    useQuery(
      trpc.organization.getHetznerCredential.queryOptions({
        id: organization.id,
      }),
    );

  const isCredentialValid = !!hetznerCredentialData?.credential?.secret;

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
        setShowDeleteDialog(false);
        navigate(`/${slug}/clusters/${id}/general`);
      },
    }),
  );

  return (
    <>
      <div className="w-full px-4">
        {isCredentialLoading ? null : isCredentialValid ? null : (
          <WarningBanner
            text="You must enter your Hetzner API Key to delete the cluster."
            linkOut={{
              text: "API Keys",
              href: `/${organization.slug}/settings/cluster/api-keys`,
            }}
            className="mt-4"
          />
        )}
        <div className="w-full pt-4 shadow-xs">
          <div className="border-mauve-6 rounded-md border text-sm">
            <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center border-b px-4 text-xs uppercase">
              Danger Zone
            </div>
            <div className="flex items-center justify-between px-4 py-2">
              <div>
                <p className="text-md font-bold">Delete this Cluster</p>
                <p className="text-mauve-11 text-xs">
                  Once you delete a cluster, there is no going back. Please be
                  certain.
                </p>
              </div>
              <Button
                className="w-36"
                intent="danger"
                disabled={
                  !isCredentialValid ||
                  clusterData?.status === "pending" ||
                  clusterData?.status === "deleted"
                }
                size="sm"
                onClick={() => setShowDeleteDialog(true)}
              >
                Delete this Cluster
              </Button>
            </div>
          </div>
        </div>
      </div>

      <DestructiveDialog
        open={showDeleteDialog}
        onOpenChange={setShowDeleteDialog}
        title="Delete this Cluster"
        description="Are you sure you want to delete this cluster? This action cannot be undone and all associated resources will be permanently removed."
        onConfirm={() => deleteClusterMutation.mutate({ id: Number(id) })}
      />
    </>
  );
}

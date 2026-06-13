import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import Button from "~/components/atoms/button/Button";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

export function ClusterAccess({ teamId }: { teamId: number }) {
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const [showAssignClusterDialog, setShowAssignClusterDialog] = useState(false);
  const [pendingAssignedClusterIds, setPendingAssignedClusterIds] = useState<
    Set<number>
  >(new Set());
  const queryClient = useQueryClient();

  const { data: teamClusters, isLoading: isTeamClustersLoading } = useQuery(
    trpc.team.getTeamClusters.queryOptions({
      teamId,
    }),
  );

  const { data: allClusters, isLoading: isAllClustersLoading } = useQuery({
    ...trpc.organization.getOrganizationClusters.queryOptions({
      id: organization.id,
    }),
    enabled: organization.isOwner,
  });

  const setTeamClustersMutation = useMutation(
    trpc.team.setTeamClusters.mutationOptions(),
  );

  const allClustersSorted =
    allClusters?.slice().sort((a, b) =>
      a.name.localeCompare(b.name, undefined, {
        sensitivity: "base",
      }),
    ) ?? [];

  function getAssignedClusterIds() {
    return new Set(teamClusters?.map((c) => c.clusterId) ?? []);
  }

  useEffect(() => {
    if (showAssignClusterDialog) {
      setPendingAssignedClusterIds(
        new Set(teamClusters?.map((c) => c.clusterId) ?? []),
      );
    }
  }, [showAssignClusterDialog, teamClusters]);

  function toggleCluster(clusterId: number, checked: boolean) {
    setPendingAssignedClusterIds((prev) => {
      const next = new Set(prev);

      if (checked) {
        next.add(clusterId);
      } else {
        next.delete(clusterId);
      }

      return next;
    });
  }

  function onApply() {
    const clusters = allClustersSorted
      .filter((cluster) => pendingAssignedClusterIds.has(cluster.id))
      .map((cluster) => ({
        clusterId: cluster.id,
      }));

    setTeamClustersMutation.mutate(
      {
        teamId,
        clusters,
      },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamClusters.queryKey(),
          });

          setShowAssignClusterDialog(false);
        },
      },
    );
  }

  function onCancel() {
    setPendingAssignedClusterIds(getAssignedClusterIds());
    setShowAssignClusterDialog(false);
  }

  return (
    <div className="rounded-md border border-mauve-6 bg-gray-2 text-sm shadow-xs">
      <div className="flex h-14 items-center justify-between rounded-t-md px-4 font-bold text-mauve-12 text-xs uppercase">
        Clusters
        {organization.isOwner && (
          <Dialog
            open={showAssignClusterDialog}
            onOpenChange={(open) => {
              setShowAssignClusterDialog(open);

              if (!open) {
                setPendingAssignedClusterIds(getAssignedClusterIds());
              }
            }}
          >
            <DialogTrigger asChild>
              <Button intent="secondary" className="w-32 text-xs">
                Manage Clusters
              </Button>
            </DialogTrigger>
            <DialogContent>
              <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                  <h1>Manage cluster access</h1>
                  <p className="text-mauve-11 text-sm">
                    Add or remove clusters to control what this team can see.
                  </p>
                </div>
                {isAllClustersLoading ? (
                  <div className="flex flex-col gap-2">
                    <Skeleton className="h-12 w-full" />
                    <Skeleton className="h-12 w-full" />
                    <Skeleton className="h-12 w-full" />
                  </div>
                ) : allClustersSorted.length === 0 ? (
                  <WarningBanner
                    text="No clusters available to be assigned."
                    linkOut={{
                      text: "Create a cluster",
                      href: `/${organization.slug}/clusters/new`,
                    }}
                  />
                ) : (
                  <div className="flex max-h-[60vh] flex-col overflow-y-auto rounded-md border border-mauve-6 bg-mauve-2">
                    {allClustersSorted.map((cluster) => (
                      <label
                        key={cluster.id}
                        className="flex min-w-0 cursor-pointer items-center gap-3 p-3"
                      >
                        <input
                          type="checkbox"
                          checked={pendingAssignedClusterIds.has(cluster.id)}
                          onChange={(event) => {
                            toggleCluster(cluster.id, event.target.checked);
                          }}
                          className="h-4.5 w-4.5 shrink-0 rounded border-mauve-6"
                        />
                        <div className="flex items-center gap-2">
                          <p
                            className="truncate font-medium text-mauve-12 text-sm"
                            title={cluster.name}
                          >
                            {cluster.name}
                          </p>
                          <span
                            className="flex items-center gap-1 text-mauve-11 text-sm"
                            title={cluster.serverType}
                          >
                            <span className="rounded-md border border-mauve-6 bg-white px-1 py-0.5">
                              {cluster.serverType}
                            </span>
                          </span>
                        </div>
                      </label>
                    ))}
                  </div>
                )}
                <div className="flex w-full justify-end gap-2">
                  <Button
                    intent="secondary"
                    className="w-24"
                    onClick={onCancel}
                    disabled={setTeamClustersMutation.isPending}
                  >
                    Cancel
                  </Button>
                  <Button
                    className="w-24"
                    onClick={onApply}
                    disabled={
                      setTeamClustersMutation.isPending ||
                      allClustersSorted.length === 0
                    }
                  >
                    Apply
                  </Button>
                </div>
              </div>
            </DialogContent>
          </Dialog>
        )}
      </div>
      <div className="mx-1 mb-1 divide-y divide-mauve-6 overflow-hidden rounded-md border border-mauve-6 bg-white shadow-xs">
        {isTeamClustersLoading ? (
          <div className="flex h-14 items-center justify-between gap-2 px-4">
            <Skeleton className="h-5 w-36" />
            <Skeleton className="h-5 w-24" />
          </div>
        ) : teamClusters?.length === 0 ? (
          <div className="flex h-14 items-center px-4 text-mauve-11 text-sm">
            No clusters assigned.
          </div>
        ) : (
          teamClusters?.map((cluster) => (
            <div
              key={cluster.clusterId}
              className="flex h-14 items-center justify-between gap-2 px-4"
            >
              <span
                className="truncate font-medium text-mauve-12"
                title={cluster.clusterName}
              >
                {cluster.clusterName}
              </span>
              <span
                className="truncate text-mauve-11"
                title={cluster.serverType}
              >
                {cluster.serverType}
              </span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

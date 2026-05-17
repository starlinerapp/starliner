import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import React, { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import { Servers } from "~/components/atoms/icons";

export function ClusterAccess({ teamId }: { teamId: number }) {
  const trpc = useTRPC();
  const organization = useOrganizationContext();
  const [showAssignClusterDialog, setShowAssignClusterDialog] = useState(false);
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

  const assignClusterMutation = useMutation(
    trpc.team.assignClusterToTeam.mutationOptions(),
  );

  const unassignClusterMutation = useMutation(
    trpc.team.unassignClusterFromTeam.mutationOptions(),
  );

  function onAssignCluster(clusterId: number) {
    assignClusterMutation.mutate(
      {
        teamId,
        clusterId,
      },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamClusters.queryKey(),
          });
        },
      },
    );
  }

  function onUnassignCluster(clusterId: number) {
    unassignClusterMutation.mutate(
      {
        teamId,
        clusterId,
      },
      {
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getTeamClusters.queryKey(),
          });
        },
      },
    );
  }

  const assignedClusterIds = new Set(
    teamClusters?.map((c) => c.clusterId) ?? [],
  );
  const allClustersSorted =
    allClusters?.slice().sort((a, b) =>
      a.name.localeCompare(b.name, undefined, {
        sensitivity: "base",
      }),
    ) ?? [];

  return (
    <div className="flex flex-col">
      <div className="w-full">
        <div className="border-mauve-6 overflow-hidden rounded-md border text-sm shadow-xs">
          <table className="w-full table-fixed border-collapse">
            <thead className="h-14">
              <tr className="border-mauve-6 bg-gray-2 border-b">
                <th className="text-mauve-12 w-[40%] px-4 py-3 text-left text-xs font-bold uppercase">
                  Cluster
                </th>
                <th className="text-mauve-12 w-[40%] px-4 py-3 text-left text-xs font-bold uppercase">
                  Server Type
                </th>
                <th className="w-[20%] px-4 py-3">
                  <div className="flex justify-end">
                    {organization.isOwner && (
                      <Button
                        intent="secondary"
                        className="w-32 text-xs"
                        onClick={() => setShowAssignClusterDialog(true)}
                      >
                        Manage Clusters
                      </Button>
                    )}
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              {isTeamClustersLoading ? (
                <tr>
                  <td className="px-4 py-3">
                    <Skeleton className="h-5 w-36" />
                  </td>

                  <td className="px-4 py-3">
                    <Skeleton className="h-5 w-24" />
                  </td>

                  <td className="px-4 py-3" />
                </tr>
              ) : teamClusters?.length === 0 ? (
                <tr>
                  <td colSpan={3} className="text-mauve-11 px-4 py-3 text-sm">
                    No clusters assigned.
                  </td>
                </tr>
              ) : (
                teamClusters?.map((cluster) => (
                  <tr
                    key={cluster.clusterId}
                    className="border-mauve-6 border-b last:border-b-0"
                  >
                    <td className="px-4 py-3">
                      <span
                        className="text-mauve-12 block truncate font-medium"
                        title={cluster.clusterName}
                      >
                        {cluster.clusterName}
                      </span>
                    </td>

                    <td className="text-mauve-11 px-4 py-3">
                      <span
                        className="block truncate"
                        title={cluster.serverType}
                      >
                        {cluster.serverType}
                      </span>
                    </td>

                    <td className="px-4 py-3" />
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
      <Dialog
        open={showAssignClusterDialog}
        onOpenChange={setShowAssignClusterDialog}
      >
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
                <Skeleton className="h-8 w-full" />
                <Skeleton className="h-8 w-full" />
                <Skeleton className="h-8 w-full" />
              </div>
            ) : allClustersSorted.length === 0 ? (
              <WarningBanner
                text={"No clusters available to be assigned."}
                linkOut={{
                  text: "Create a cluster",
                  href: `/${organization.slug}/clusters/new`,
                }}
              />
            ) : (
              <div className="flex max-h-[60vh] flex-col gap-1 overflow-y-auto">
                {allClustersSorted.map((cluster) => {
                  const isAssigned = assignedClusterIds.has(cluster.id);
                  return (
                    <div
                      key={cluster.id}
                      className="bg-mauve-2 border-mauve-6 flex min-w-0 items-center justify-between gap-3 rounded-md border p-3"
                    >
                      <div className="border-mauve-6 rounded-md border bg-white p-1.5">
                        <Servers className="text-mauve-11 h-7 w-7" />
                      </div>
                      <div className="flex min-w-0 flex-1 flex-col">
                        <p
                          className="text-mauve-12 truncate text-sm font-medium"
                          title={cluster.name}
                        >
                          {cluster.name}
                        </p>
                        <span
                          className="text-mauve-11 flex items-center gap-1 text-xs"
                          title={cluster.serverType}
                        >
                          <p>Server Type:</p>
                          <span className="border-mauve-6 rounded-md border bg-white px-1 py-0.5">
                            {cluster.serverType}
                          </span>
                        </span>
                      </div>
                      {isAssigned ? (
                        <Button
                          className="w-24"
                          size="xs"
                          intent="secondary"
                          onClick={() => onUnassignCluster(cluster.id)}
                        >
                          Unassign
                        </Button>
                      ) : (
                        <Button
                          className="w-24"
                          size="xs"
                          intent="primary"
                          onClick={() => onAssignCluster(cluster.id)}
                        >
                          Assign
                        </Button>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}

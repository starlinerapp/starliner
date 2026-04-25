import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { Cross } from "~/components/atoms/icons";
import React, { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

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
  const unassignedClusters =
    allClusters?.filter((c) => !assignedClusterIds.has(c.id)) ?? [];
  return (
    <div className="w-full">
      <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
          <p>Cluster Access</p>
          <Dialog
            open={showAssignClusterDialog}
            onOpenChange={setShowAssignClusterDialog}
          >
            <DialogTrigger asChild>
              <Button intent="secondary" className="h-7 w-28 text-xs">
                Assign Cluster
              </Button>
            </DialogTrigger>
            <DialogContent>
              <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                  <h1>Assign Cluster</h1>
                  <p className="text-mauve-11 text-sm">
                    Select a cluster to make visible to this team&apos;s
                    members.
                  </p>
                </div>
                {isAllClustersLoading ? (
                  <div className="flex flex-col gap-2">
                    <Skeleton className="h-8 w-full" />
                    <Skeleton className="h-8 w-full" />
                    <Skeleton className="h-8 w-full" />
                  </div>
                ) : allClusters?.length === 0 ? (
                  <div className="text-mauve-11 text-sm">
                    No clusters exist yet. Create a cluster first before
                    assigning it to this team.
                  </div>
                ) : unassignedClusters.length === 0 ? (
                  <div className="text-mauve-11 text-sm">
                    All clusters are already assigned to this team.
                  </div>
                ) : (
                  <div className="border-mauve-6 divide-mauve-6 max-h-[60vh] divide-y overflow-y-auto rounded-md border">
                    {unassignedClusters.map((cluster) => (
                      <div
                        key={cluster.id}
                        className="flex items-center justify-between px-2 py-2"
                      >
                        <div className="flex flex-col">
                          <span className="text-mauve-12 text-sm font-medium">
                            {cluster.name}
                          </span>
                          <span className="text-mauve-11 text-sm">
                            {cluster.serverType}
                          </span>
                        </div>
                        <Button
                          className="h-7 w-20 text-xs"
                          intent="secondary"
                          disabled={assignClusterMutation.isPending}
                          onClick={() => {
                            onAssignCluster(cluster.id);
                          }}
                        >
                          Assign
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </DialogContent>
          </Dialog>
        </div>
        {isTeamClustersLoading ? (
          <div className="flex flex-col gap-2 px-4 py-3">
            <Skeleton className="h-5 w-48" />
            <Skeleton className="h-5 w-36" />
          </div>
        ) : teamClusters?.length === 0 ? (
          <div className="text-mauve-11 px-4 py-3 text-sm">
            No clusters assigned. Team members cannot see any clusters until you
            assign them.
          </div>
        ) : (
          teamClusters?.map((cluster) => (
            <div
              key={cluster.clusterId}
              className="border-mauve-6 text-mauve-12 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
            >
              <div className="flex flex-col">
                <span>{cluster.clusterName}</span>
                <span className="text-mauve-11 text-sm">
                  {cluster.serverType}
                </span>
              </div>
              <button
                className="text-mauve-11 hover:text-mauve-12 cursor-pointer"
                disabled={unassignClusterMutation.isPending}
                onClick={() => onUnassignCluster(cluster.clusterId)}
                title="Revoke cluster access"
              >
                <Cross width={16} height={16} />
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

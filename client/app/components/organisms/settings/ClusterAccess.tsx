import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import React, { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import WarningBanner from "~/components/atoms/banner/WarningBanner";

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
    allClusters
      ?.slice()
      .sort((a, b) =>
        a.name.localeCompare(b.name, undefined, { sensitivity: "base" }),
      ) ?? [];

  return (
    <div className="w-full">
      <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center justify-between border-b px-4 text-xs font-bold uppercase">
          <p>Cluster Access</p>
          <Dialog
            open={showAssignClusterDialog}
            onOpenChange={setShowAssignClusterDialog}
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
                  <div className="border-mauve-6 divide-mauve-6 max-h-[60vh] divide-y overflow-y-auto rounded-md border">
                    {allClustersSorted.map((cluster) => {
                      const isAssigned = assignedClusterIds.has(cluster.id);
                      return (
                        <div
                          key={cluster.id}
                          className="flex min-w-0 items-center justify-between gap-3 p-3"
                        >
                          <div className="flex min-w-0 flex-1 flex-col gap-1">
                            <span
                              className="text-mauve-12 truncate text-sm font-medium"
                              title={cluster.name}
                            >
                              {cluster.name}
                            </span>
                            <span
                              className="text-mauve-11 truncate text-xs"
                              title={cluster.serverType}
                            >
                              Server Type: {cluster.serverType}
                            </span>
                          </div>
                          {isAssigned ? (
                            <Button
                              className="w-24 text-xs"
                              intent="secondary"
                              onClick={() => onUnassignCluster(cluster.id)}
                            >
                              Unassign
                            </Button>
                          ) : (
                            <Button
                              className="w-24 text-xs"
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
        {isTeamClustersLoading ? (
          <div className="flex flex-col gap-1 px-4 py-3">
            <Skeleton className="h-5 w-48" />
            <Skeleton className="h-5 w-36" />
          </div>
        ) : teamClusters?.length === 0 ? (
          <p className="text-mauve-11 px-4 py-3 text-sm">
            No clusters assigned. <br /> Team members cannot see any clusters
            until you assign them.
          </p>
        ) : (
          teamClusters?.map((cluster) => (
            <div
              key={cluster.clusterId}
              className="border-mauve-6 text-mauve-12 min-w-0 border-b px-4 py-3 text-sm last:border-b-0"
            >
              <div className="flex min-w-0 flex-col">
                <span className="truncate" title={cluster.clusterName}>
                  {cluster.clusterName}
                </span>
                <span
                  className="text-mauve-11 truncate text-sm"
                  title={cluster.serverType}
                >
                  {cluster.serverType}
                </span>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

import React from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

export default function ProjectSettings() {
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const { slug, id } = useParams<{
    slug: string;
    id: string;
  }>();

  const deleteProjectMutation = useMutation(
    trpc.project.deleteProject.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${slug}/projects/all`);
      },
    }),
  );

  const { data: clusterData, isLoading: isProjectDataLoading } = useQuery(
    trpc.project.getProjectCluster.queryOptions({ id: Number(id) }),
  );

  return (
    <div className="w-full space-y-4 p-4 xl:w-3/5">
      <div className="border-mauve-6 rounded-md border-1 text-sm">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs uppercase">
          General
        </div>
        <div className="flex items-center justify-between px-4 py-2">
          <div>Assigned Cluster</div>
          {isProjectDataLoading ? (
            <Skeleton className="h-10 w-60" />
          ) : (
            <input
              className="border-mauve-6 w-60 cursor-not-allowed rounded-md border-1 p-2"
              placeholder={clusterData.clusterName}
              disabled
            />
          )}
        </div>
      </div>
      <div className="border-mauve-6 rounded-md border-1 text-sm">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs uppercase">
          Danger Zone
        </div>
        <div className="flex items-center justify-between px-4 py-2">
          <div>
            <p className="text-md font-bold">Delete this Project</p>
            <p className="text-mauve-11 text-xs">
              Once you delete a project, there is no going back. Please be
              certain.
            </p>
          </div>
          <Button
            className="w-36"
            intent="danger"
            size="sm"
            onClick={() =>
              deleteProjectMutation.mutate({
                id: Number(id),
              })
            }
          >
            Delete this Project
          </Button>
        </div>
      </div>
    </div>
  );
}

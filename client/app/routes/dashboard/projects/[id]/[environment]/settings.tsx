import React, { useEffect, useMemo, useState } from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import UpdateConnectedBranchForm from "~/components/organisms/forms/UpdateConnectedBranchForm";
import Switch from "~/components/atoms/switch/Switch";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";

export default function ProjectSettings() {
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const {
    slug,
    id,
    environment: environmentSlug,
  } = useParams<{
    slug: string;
    id: string;
    environment: string;
  }>();

  const projectId = Number(id);

  const { data: project } = useQuery(
    trpc.project.getProject.queryOptions({ id: projectId }),
  );
  const environmentId = useMemo(
    () => project?.environments.find((e) => e.slug === environmentSlug)?.id,
    [project, environmentSlug],
  );

  const deleteProjectMutation = useMutation(
    trpc.project.deleteProject.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getUserProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${slug}/projects/all`);
      },
    }),
  );

  const deleteEnvironmentMutation = useMutation(
    trpc.environment.deleteEnvironment.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getUserProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${slug}/projects/${id}`);
      },
    }),
  );

  const { data: clusterData, isLoading: isClusterDataLoading } = useQuery(
    trpc.project.getProjectCluster.queryOptions({ id: projectId }),
  );

  const { data: previewEnvEnabled, isLoading: isPreviewEnvEnabledLoading } =
    useQuery(
      trpc.project.getProjectPreviewEnvironmentEnabled.queryOptions({
        id: projectId,
      }),
    );

  const [showPreviewEnvDialog, setShowPreviewEnvDialog] = useState(false);
  const [pendingCheckedValue, setPendingCheckedValue] = useState<
    boolean | null
  >(null);

  useEffect(() => {
    if (!showPreviewEnvDialog) {
      setPendingCheckedValue(null);
    }
  }, [showPreviewEnvDialog]);

  const togglePreviewEnvMutation = useMutation(
    trpc.project.toggleProjectPreviewEnvironmentEnabled.mutationOptions({
      onSuccess: async () => {
        setShowPreviewEnvDialog(false);
        setPendingCheckedValue(null);

        await queryClient.invalidateQueries({
          queryKey: trpc.project.getProjectPreviewEnvironmentEnabled.queryKey({
            id: projectId,
          }),
        });
      },
      onError: () => {
        setShowPreviewEnvDialog(false);
        setPendingCheckedValue(null);
      },
    }),
  );

  return (
    <>
      <div className="w-full space-y-4 p-4">
        <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center border-b px-4 text-xs uppercase">
            Project Settings
          </div>
          <div className="border-mauve-6 flex items-center justify-between px-4 py-2">
            <div className="flex flex-col">
              <p className="text-md font-bold">PR Environments</p>
              <p className="text-mauve-11 text-xs">
                Automatically created by Starliner when a pull request is opened
                and up when the PR is closed.
              </p>
            </div>

            {isPreviewEnvEnabledLoading ? (
              <Skeleton className="h-6.25 w-10.5 rounded-full" />
            ) : (
              <Switch
                checked={previewEnvEnabled?.enabled ?? false}
                disabled={togglePreviewEnvMutation.isPending}
                onCheckedChange={(checked) => {
                  setPendingCheckedValue(checked);
                  setShowPreviewEnvDialog(true);
                }}
              />
            )}
          </div>
        </div>

        <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center border-b px-4 text-xs uppercase">
            Environment Settings
          </div>
          <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
            <div className="flex flex-col">
              <p className="text-md font-bold">Assigned Cluster</p>
              <p className="text-mauve-11 text-xs">
                The Cluster this project is running on.
              </p>
            </div>
            {isClusterDataLoading ? (
              <Skeleton className="h-9.5 w-1/2" />
            ) : (
              <input
                className="border-mauve-6 disabled:text-mauve-11 w-1/2 cursor-not-allowed rounded-md border p-2"
                value={clusterData?.clusterName}
                disabled
              />
            )}
          </div>
          <UpdateConnectedBranchForm />
        </div>

        <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center border-b px-4 text-xs uppercase">
            Danger Zone
          </div>
          <div className="border-mauve-6 flex items-center justify-between border-b px-4 py-2">
            <div>
              <p className="text-md font-bold">Delete this Environment</p>
              <p className="text-mauve-11 text-xs">
                Once you delete an environment, there is no going back. Please
                be certain.
              </p>
            </div>
            <Button
              className="w-46"
              intent="danger"
              disabled={
                deleteEnvironmentMutation.isPending ||
                environmentId == null ||
                isClusterDataLoading ||
                environmentSlug === "production"
              }
              size="sm"
              onClick={() => {
                if (environmentId == null) return;
                deleteEnvironmentMutation.mutate({
                  id: environmentId,
                });
              }}
            >
              Delete this Environment
            </Button>
          </div>
          <div className="border-mauve-6 flex items-center justify-between px-4 py-2">
            <div>
              <p className="text-md font-bold">Delete this Project</p>
              <p className="text-mauve-11 text-xs">
                Deleting the project will delete all environments and
                deployments associated with it.
              </p>
            </div>
            <Button
              className="w-36"
              intent="danger"
              disabled={isClusterDataLoading}
              size="sm"
              onClick={() =>
                deleteProjectMutation.mutate({
                  id: projectId,
                })
              }
            >
              Delete this Project
            </Button>
          </div>
        </div>
      </div>

      <Dialog
        open={showPreviewEnvDialog}
        onOpenChange={(open) => {
          setShowPreviewEnvDialog(open);
          if (!open) {
            setPendingCheckedValue(null);
          }
        }}
      >
        <DialogContent>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <h1>
                {pendingCheckedValue
                  ? "Enable PR Environments"
                  : "Disable PR Environments"}
              </h1>
              <p className="text-mauve-11 text-sm">
                {pendingCheckedValue
                  ? "Are you sure you want to enable preview environments for this project?"
                  : "Are you sure you want to disable preview environments for this project?"}
              </p>
            </div>

            <div className="flex justify-end gap-2">
              <Button
                type="button"
                intent="secondary"
                className="w-24"
                onClick={() => {
                  setShowPreviewEnvDialog(false);
                  setPendingCheckedValue(null);
                }}
              >
                Cancel
              </Button>
              <Button
                className="w-24"
                disabled={togglePreviewEnvMutation.isPending}
                onClick={() => {
                  togglePreviewEnvMutation.mutate({
                    id: projectId,
                  });
                }}
              >
                Confirm
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

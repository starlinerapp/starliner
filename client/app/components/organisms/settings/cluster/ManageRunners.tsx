import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import Button from "~/components/atoms/button/Button";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import NewRunnerDialog from "~/components/organisms/settings/cluster/NewRunnerDialog";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

function RunnerStatus({ status }: { status: string }) {
  const isOnline = status === "online";

  return (
    <span className="flex items-center gap-1.5 text-mauve-11 text-sm capitalize">
      <span
        className={cn(
          "size-2 rounded-full",
          isOnline ? "bg-grass-9" : "bg-mauve-8",
        )}
      />
      {status}
    </span>
  );
}

export default function ManageRunners() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const organization = useOrganizationContext();
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  const { data: runnersData, isLoading: isRunnersLoading } = useQuery(
    trpc.runner.getOrganizationRunners.queryOptions({
      organizationId: organization.id,
    }),
  );

  const registeredRunners = runnersData?.filter((runner) => runner.name) ?? [];

  function handleCreateDialogChange(open: boolean) {
    setShowCreateDialog(open);
    if (!open) {
      void queryClient.invalidateQueries({
        queryKey: trpc.runner.getOrganizationRunners.queryKey({
          organizationId: organization.id,
        }),
      });
    }
  }

  return (
    <div className="flex flex-col">
      <div className="rounded-md border border-mauve-6 bg-gray-2 text-sm shadow-xs">
        <div className="flex h-14 items-center justify-between rounded-t-md px-4 font-bold text-mauve-12 text-xs uppercase">
          Runners
          {organization.isOwner && (
            <Button
              className="text-xs"
              intent="secondary"
              onClick={() => setShowCreateDialog(true)}
            >
              New self-hosted runner
            </Button>
          )}
        </div>
        <div className="mx-1 mb-1 divide-y divide-mauve-6 overflow-hidden rounded-md border border-mauve-6 bg-white shadow-xs">
          {isRunnersLoading ? (
            Array.from({ length: 2 }).map((_, i) => (
              <div key={i} className="flex h-14 items-center gap-3 px-4">
                <Skeleton className="h-9 w-9 rounded-md" />
                <div className="flex flex-1 flex-col gap-1.5">
                  <Skeleton className="h-3.5 w-32" />
                  <Skeleton className="h-3 w-48" />
                </div>
              </div>
            ))
          ) : registeredRunners.length === 0 ? (
            <div className="flex h-14 items-center px-4 text-mauve-11 text-sm">
              No runners yet.
            </div>
          ) : (
            registeredRunners.map((runner) => (
              <div
                key={runner.id}
                className="flex h-14 items-center gap-3 px-4"
              >
                <div className="flex h-9 w-9 items-center justify-center rounded-md bg-violet-9 text-base text-white">
                  {runner.name?.substring(0, 1)?.toUpperCase()}
                </div>
                <div className="flex min-w-0 flex-1 flex-col gap-0.5">
                  <span className="truncate font-medium text-mauve-12 text-sm">
                    {runner.name}
                  </span>
                  {runner.labels.length > 0 && (
                    <div className="flex max-w-full flex-wrap gap-0.5 self-start">
                      {runner.labels.map((label) => (
                        <span
                          key={label}
                          className="max-w-40 truncate rounded-md border border-mauve-6 px-1 text-mauve-11 text-xs"
                        >
                          {label}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
                <RunnerStatus status={runner.status} />
              </div>
            ))
          )}
        </div>
      </div>
      {organization.isOwner && (
        <Dialog open={showCreateDialog} onOpenChange={handleCreateDialogChange}>
          <DialogContent>
            <NewRunnerDialog />
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

import { useState } from "react";
import Button from "~/components/atoms/button/Button";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import NewRunnerDialog from "~/components/organisms/settings/cluster/NewRunnerDialog";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

export default function ManageRunners() {
  const organization = useOrganizationContext();
  const [showCreateDialog, setShowCreateDialog] = useState(false);

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
          <div className="flex h-14 items-center px-4 text-mauve-11 text-sm">
            No runners yet.
          </div>
        </div>
      </div>
      {organization.isOwner && (
        <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
          <DialogContent>
            <NewRunnerDialog />
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

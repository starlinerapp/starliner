import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useNavigate } from "react-router";
import Button from "~/components/atoms/button/Button";
import DestructiveDialog from "~/components/organisms/dialog/DestructiveDialog";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

interface TeamDangerZoneProps {
  teamId: number;
}

export default function TeamDangerZone({ teamId }: TeamDangerZoneProps) {
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const deleteTeamMutation = useMutation(
    trpc.team.deleteTeam.mutationOptions({
      onSuccess: async () => {
        setShowDeleteDialog(false);
        await queryClient.invalidateQueries({
          queryKey: trpc.team.getUserTeams.queryKey({
            organizationId: organization.id,
          }),
        });
        navigate("../", { relative: "path" });
      },
    }),
  );

  function handleDeleteTeam() {
    deleteTeamMutation.reset();
    setShowDeleteDialog(true);
  }

  function confirmDeleteTeam() {
    deleteTeamMutation.mutate({ teamId });
  }

  return (
    <>
      <div className="rounded-md border border-mauve-6 bg-gray-2 text-sm shadow-xs">
        <div className="flex h-14 items-center rounded-t-md px-4 font-bold text-mauve-12 text-xs uppercase">
          Danger Zone
        </div>
        <div className="mx-1 mb-1 overflow-hidden rounded-md border border-mauve-6 bg-white shadow-xs">
          <div className="flex h-14 items-center justify-between gap-2 px-4">
            <div>
              <h2 className="text-mauve-12">Delete this Team</h2>
              <p className="text-mauve-11 text-xs">
                Once you delete a team, there is no going back. Please be
                certain.
              </p>
            </div>
            <Button
              className="w-36"
              intent="danger"
              size="sm"
              disabled={deleteTeamMutation.isPending}
              onClick={handleDeleteTeam}
            >
              Delete this Team
            </Button>
          </div>
        </div>
      </div>

      <DestructiveDialog
        open={showDeleteDialog}
        onOpenChange={(open) => {
          setShowDeleteDialog(open);
          if (!open) {
            deleteTeamMutation.reset();
          }
        }}
        title={"Delete Team"}
        bannerText={
          "Deleting the team will delete all projects that belong to the team, including the deployments."
        }
        description={
          "Are you sure you want to delete this team? This action cannot be undone."
        }
        isPending={deleteTeamMutation.isPending}
        onConfirm={confirmDeleteTeam}
      />
    </>
  );
}

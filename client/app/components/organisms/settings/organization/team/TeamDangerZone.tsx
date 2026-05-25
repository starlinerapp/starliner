import React, { useState } from "react";
import Button from "~/components/atoms/button/Button";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import DestructiveDialog from "~/components/organisms/dialog/DestructiveDialog";

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
      <div className="border-mauve-6 overflow-hidden rounded-md border text-sm shadow-xs">
        <table className="w-full table-fixed border-collapse">
          <thead className="h-14">
            <tr className="border-mauve-6 bg-gray-2 border-b">
              <th className="text-mauve-12 w-[40%] px-4 py-3 text-left text-xs font-bold uppercase">
                Danger Zone
              </th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="flex justify-between px-4 py-2">
                <span>
                  <p className="text-md font-bold">Delete this Team</p>
                  <p className="text-mauve-11 text-xs">
                    Once you delete a team, there is no going back. Please be
                    certain.
                  </p>
                </span>
                <Button
                  className="w-36"
                  intent="danger"
                  size="sm"
                  disabled={deleteTeamMutation.isPending}
                  onClick={handleDeleteTeam}
                >
                  Delete this Team
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
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

import React, { useState } from "react";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { Link, useParams } from "react-router";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import { formatSlugInput } from "~/utils/slug";

interface CreateTeamFormInput {
  name: string;
}

export default function OrganizationTeams() {
  const trpc = useTRPC();
  const { slug } = useParams();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();
  const [showCreateDialog, setShowCreateDialog] = useState(false);

  const {
    register: registerCreate,
    handleSubmit: handleCreateSubmit,
    reset: resetCreate,
    watch: watchCreate,
    setValue: setCreateValue,
  } = useForm<CreateTeamFormInput>();
  const nameInput = watchCreate("name", "");

  const { data: teamsData, isLoading } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );

  const createTeamMutation = useMutation(
    trpc.team.createTeam.mutationOptions(),
  );

  function onCreateTeam(data: CreateTeamFormInput) {
    createTeamMutation.mutate(
      {
        organizationId: organization.id,
        name: data.name,
      },
      {
        onSuccess: async () => {
          resetCreate();
          setShowCreateDialog(false);
          await queryClient.invalidateQueries({
            queryKey: trpc.team.getUserTeams.queryKey(),
          });
        },
      },
    );
  }

  return (
    <div className="w-full xl:w-3/5">
      <div className="border-mauve-6 w-full rounded-md border-1 text-sm">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
          <p>Your Teams</p>
          <div className="flex gap-2">
            <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
              <DialogTrigger asChild>
                <Button className="h-7 w-24 text-xs">Create Team</Button>
              </DialogTrigger>
              <DialogContent>
                <h2 className="text-mauve-12 mb-4 text-lg font-bold">
                  Create Team
                </h2>
                <form
                  className="flex flex-col gap-3"
                  onSubmit={handleCreateSubmit(onCreateTeam)}
                >
                  <input
                    className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
                    placeholder="team-slug"
                    maxLength={50}
                    {...registerCreate("name")}
                    onChange={(e) => {
                      setCreateValue("name", formatSlugInput(e.target.value));
                    }}
                  />
                  <div className="flex justify-end gap-2">
                    <Button
                      className="h-9 w-24 text-xs"
                      intent="secondary"
                      type="button"
                      onClick={() => {
                        setShowCreateDialog(false);
                        resetCreate();
                      }}
                    >
                      Cancel
                    </Button>
                    <Button
                      className="h-9 w-24 text-xs"
                      type="submit"
                      disabled={!nameInput || createTeamMutation.isPending}
                    >
                      Save
                    </Button>
                  </div>
                </form>
              </DialogContent>
            </Dialog>
          </div>
        </div>
        {isLoading ? (
          <div className="text-mauve-11 px-4 py-3 text-sm">Loading...</div>
        ) : teamsData?.length === 0 ? (
          <div className="text-mauve-11 px-4 py-3 text-sm">No teams yet.</div>
        ) : (
          teamsData?.map((team) => (
            <Link
              key={team.id}
              to={`/${slug}/settings/teams/${team.id}`}
              className="border-mauve-6 text-mauve-12 hover:bg-gray-2 flex items-center justify-between border-b px-4 py-3 text-sm last:border-b-0"
            >
              <div className="flex flex-col">
                <span>#{team.slug}</span>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}

import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Button from "~/components/atoms/button/Button";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
} from "~/components/atoms/dialog/Dialog";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { formatSlugInput, sanitizeSlug } from "~/utils/slug";
import OrganizationTeams from "~/components/organisms/settings/OrganizationTeams";

interface CreateTeamFormInput {
  name: string;
}

export default function Teams() {
  const trpc = useTRPC();
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

  const createTeamMutation = useMutation(
    trpc.team.createTeam.mutationOptions(),
  );

  function onCreateTeam(data: CreateTeamFormInput) {
    createTeamMutation.mutate(
      {
        organizationId: organization.id,
        name: sanitizeSlug(data.name),
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
    <div className="flex w-full flex-col px-8 py-4">
      <div className="flex min-h-10 w-full items-center justify-between">
        <h1 className="text-xl font-bold">Teams</h1>
        {organization.isOwner && (
          <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
            <DialogTrigger asChild>
              <Button className="w-32">Create Team</Button>
            </DialogTrigger>
            <DialogContent>
              <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                  <h1>Create Team</h1>
                  <p className="text-mauve-11 text-sm">
                    Teams group members for project and cluster visibility.
                  </p>
                </div>
                {createTeamMutation.isError && (
                  <ErrorBanner text={createTeamMutation.error.message} />
                )}
                <form
                  className="flex flex-col gap-3"
                  onSubmit={handleCreateSubmit(onCreateTeam)}
                >
                  <input
                    className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
                    placeholder="Team Slug*"
                    maxLength={50}
                    {...registerCreate("name")}
                    onChange={(e) => {
                      setCreateValue("name", formatSlugInput(e.target.value));
                      createTeamMutation.reset();
                    }}
                  />
                  <div className="flex justify-end gap-2">
                    <Button
                      intent="secondary"
                      className="w-24"
                      type="button"
                      onClick={() => {
                        setShowCreateDialog(false);
                        resetCreate();
                        createTeamMutation.reset();
                      }}
                    >
                      Cancel
                    </Button>
                    <Button
                      intent="primary"
                      className="w-24"
                      type="submit"
                      disabled={!nameInput || createTeamMutation.isPending}
                    >
                      Create
                    </Button>
                  </div>
                </form>
              </div>
            </DialogContent>
          </Dialog>
        )}
      </div>
      <div className="flex flex-col gap-4 pt-10.5">
        <OrganizationTeams />
      </div>
    </div>
  );
}

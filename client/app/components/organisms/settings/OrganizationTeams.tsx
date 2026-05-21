import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useNavigate, useParams } from "react-router";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import Button from "~/components/atoms/button/Button";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { formatSlugInput, sanitizeSlug } from "~/utils/slug";

interface CreateTeamFormInput {
  name: string;
}

export default function OrganizationTeams() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const { slug } = useParams();
  const navigate = useNavigate();
  const organization = useOrganizationContext();
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

  const { data: teamsData, isLoading: isTeamsLoading } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );

  const { data: members } = useQuery(
    trpc.organization.getOrganizationMembers.queryOptions({
      id: organization.id,
    }),
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

  const columnCount = organization.isOwner ? 3 : 2;

  return (
    <div className="flex flex-col">
      <div className="border-mauve-6 overflow-hidden rounded-md border text-sm shadow-xs">
        <table className="w-full border-collapse">
          <thead className="h-14">
            <tr className="border-mauve-6 bg-gray-2 border-b">
              <th className="text-mauve-12 w-1/2 px-4 py-3 text-left text-xs font-bold uppercase">
                Teams
              </th>
              <th className="w-1/2 px-4 py-3"></th>
              {organization.isOwner && (
                <th className="w-[20%] px-4">
                  <Button
                    className="w-28 text-xs"
                    intent="secondary"
                    onClick={() => setShowCreateDialog(true)}
                  >
                    Create Team
                  </Button>
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {isTeamsLoading ? (
              Array.from({ length: 2 }).map((_, i) => (
                <tr key={i} className="border-mauve-6 border-b last:border-b-0">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <Skeleton className="h-9 w-9 rounded-md" />
                      <span className="flex flex-col gap-1">
                        <Skeleton className="h-3.5 w-24" />
                        <Skeleton className="h-3.5 w-20" />
                      </span>
                    </div>
                  </td>
                  <td className="px-4 py-3" />
                  {organization.isOwner && <td className="px-4 py-3" />}
                </tr>
              ))
            ) : teamsData?.length === 0 ? (
              <tr>
                <td
                  colSpan={columnCount}
                  className="text-mauve-11 px-4 py-3 text-sm"
                >
                  No teams yet.
                </td>
              </tr>
            ) : (
              teamsData?.map((team) => (
                <tr
                  key={team.id}
                  className="border-mauve-6 hover:bg-gray-2 cursor-pointer border-b last:border-b-0"
                  onClick={() =>
                    navigate(`/${slug}/settings/organization/teams/${team.id}`)
                  }
                >
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <div className="bg-violet-9 flex h-9 w-9 items-center justify-center rounded-md text-base text-white">
                        {team.slug.substring(0, 1)?.toUpperCase()}
                      </div>
                      <div className="flex flex-col">
                        <span className="text-violet-11">#{team.slug}</span>
                        <p className="text-mauve-11 text-sm">
                          {members?.length ?? 0}{" "}
                          {(members?.length ?? 0) === 1 ? "Member" : "Members"}
                        </p>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3" />
                  {organization.isOwner && <td />}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      {organization.isOwner && (
        <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
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
  );
}

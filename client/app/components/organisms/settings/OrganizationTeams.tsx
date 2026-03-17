import React, { useState } from "react";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { Link, useParams } from "react-router";

interface CreateTeamFormInput {
  name: string;
}

interface JoinTeamFormInput {
  slug: string;
}

export default function OrganizationTeams() {
  const trpc = useTRPC();
  const { slug } = useParams();
  const organization = useOrganizationContext();
  const queryClient = useQueryClient();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showJoinForm, setShowJoinForm] = useState(false);

  const {
    register: registerCreate,
    handleSubmit: handleCreateSubmit,
    reset: resetCreate,
    watch: watchCreate,
  } = useForm<CreateTeamFormInput>();
  const nameInput = watchCreate("name", "");

  const {
    register: registerJoin,
    handleSubmit: handleJoinSubmit,
    reset: resetJoin,
    watch: watchJoin,
  } = useForm<JoinTeamFormInput>();
  const slugInput = watchJoin("slug", "");

  const { data: teamsData, isLoading } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );

  const createTeamMutation = useMutation(
    trpc.team.createTeam.mutationOptions(),
  );

  const joinTeamMutation = useMutation(trpc.team.joinTeam.mutationOptions());

  const leaveTeamMutation = useMutation(
    trpc.team.removeTeamMember.mutationOptions(),
  );

  return (
    <div className="w-full xl:w-3/5">
      <div className="border-mauve-6 w-full rounded-md border-1 text-sm">
        <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex items-center justify-between border-b px-4 py-2 text-xs font-bold uppercase">
          <p>Your Teams</p>
          <div className="flex gap-2">
            <Button
              className="h-7 w-24 text-xs"
              intent="secondary"
              onClick={() => {
                setShowJoinForm(!showJoinForm);
                setShowCreateForm(false);
              }}
            >
              Join Team
            </Button>
            <Button
              className="h-7 w-24 text-xs"
              onClick={() => {
                setShowCreateForm(!showCreateForm);
                setShowJoinForm(false);
              }}
            >
              Create Team
            </Button>
          </div>
        </div>
        {showCreateForm && (
          <form
            className="flex items-center gap-2 px-4 py-3"
            onSubmit={handleCreateSubmit((data) => {
              createTeamMutation.mutate(
                {
                  organizationId: organization.id,
                  name: data.name,
                },
                {
                  onSuccess: async () => {
                    resetCreate();
                    setShowCreateForm(false);
                    await queryClient.invalidateQueries({
                      queryKey: trpc.team.getUserTeams.queryKey(),
                    });
                  },
                },
              );
            })}
          >
            <input
              className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
              placeholder="Team name"
              {...registerCreate("name")}
            />
            <Button
              className="h-9 w-24 text-xs"
              type="submit"
              disabled={!nameInput || createTeamMutation.isPending}
            >
              Save
            </Button>
            <Button
              className="h-9 w-24 text-xs"
              intent="secondary"
              onClick={() => {
                setShowCreateForm(false);
                resetCreate();
              }}
            >
              Cancel
            </Button>
          </form>
        )}
        {showJoinForm && (
          <form
            className="flex items-center gap-2 px-4 py-3"
            onSubmit={handleJoinSubmit((data) => {
              joinTeamMutation.mutate(
                {
                  organizationId: organization.id,
                  slug: data.slug,
                },
                {
                  onSuccess: async () => {
                    resetJoin();
                    setShowJoinForm(false);
                    await queryClient.invalidateQueries({
                      queryKey: trpc.team.getUserTeams.queryKey(),
                    });
                  },
                },
              );
            })}
          >
            <input
              className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full rounded-md border p-2 text-sm"
              placeholder="Team slug"
              {...registerJoin("slug")}
            />
            <Button
              className="h-9 w-24 text-xs"
              type="submit"
              disabled={!slugInput || joinTeamMutation.isPending}
            >
              Join
            </Button>
            <Button
              className="h-9 w-24 text-xs"
              intent="secondary"
              onClick={() => {
                setShowJoinForm(false);
                resetJoin();
              }}
            >
              Cancel
            </Button>
          </form>
        )}
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
              <span>{team.name}</span>
              <span className="text-mauve-11 text-xs">{team.slug}</span>
              <Button
                className="h-7 w-24 text-xs"
                intent="secondary"
                disabled={leaveTeamMutation.isPending}
                onClick={(e) => {
                  e.preventDefault();
                  leaveTeamMutation.mutate(
                    {
                      organizationId: organization.id,
                      teamId: team.id,
                    },
                    {
                      onSuccess: async () => {
                        await queryClient.invalidateQueries({
                          queryKey: trpc.team.getUserTeams.queryKey(),
                        });
                      },
                    },
                  );
                }}
              >
                Leave Team
              </Button>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}

import React from "react";
import Button from "~/components/atoms/button/Button";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useNavigate } from "react-router";
import { ChevronDown } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";
import WarningBanner from "~/components/atoms/banner/WarningBanner";

interface NewProjectFormInput {
  name: string;
  clusterId: string;
  teamId: string;
}

export default function NewProject() {
  const queryClient = useQueryClient();
  const trpc = useTRPC();
  const navigate = useNavigate();

  const organization = useOrganizationContext();

  const { register, handleSubmit, watch, setValue } =
    useForm<NewProjectFormInput>();
  const nameInput = watch("name", "");
  const clusterIdInput = watch("clusterId", "");
  const teamIdInput = watch("teamId", "");

  const { data: teamsData, isLoading: isTeamsLoading } = useQuery(
    trpc.team.getUserTeams.queryOptions({
      organizationId: organization.id,
    }),
  );

  const { data: clustersData, isLoading: isClustersLoading } = useQuery(
    trpc.team.getTeamClusters.queryOptions(
      { teamId: Number(teamIdInput) },
      { enabled: !!teamIdInput },
    ),
  );

  React.useEffect(() => {
    if (teamsData?.[0] && !teamIdInput) {
      setValue("teamId", String(teamsData[0].id));
    }
  }, [teamsData, teamIdInput, setValue]);

  React.useEffect(() => {
    setValue("clusterId", "");
  }, [teamIdInput, setValue]);

  const isLoading = isClustersLoading || isTeamsLoading;

  const createProjectMutation = useMutation(
    trpc.project.createProject.mutationOptions({
      onSuccess: async (project) => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getUserProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${organization.slug}/projects/${project.id}`);
      },
    }),
  );

  const onSubmit: SubmitHandler<NewProjectFormInput> = (data) => {
    createProjectMutation.mutate({
      organizationId: organization.id,
      name: data.name,
      clusterId: Number(data.clusterId),
      teamId: Number(data.teamId),
    });
  };

  const clusterExists = clustersData && clustersData.length > 0;
  const teamExists = teamsData && teamsData.length > 0;

  return (
    <div className="flex flex-col gap-2 px-8 py-4">
      <h1 className="text-xl font-bold">New Project</h1>
      {!teamExists && !isLoading ? (
        <WarningBanner
          text="You must join or create a team before creating projects."
          linkOut={{
            text: "Manage Teams",
            href: `/${organization.slug}/settings/organization`,
          }}
          className="my-2"
        />
      ) : null}
      {teamExists && !clusterExists && !isLoading && teamIdInput ? (
        <WarningBanner
          text="This team has no clusters assigned. Assign one before creating projects."
          linkOut={{
            text: "Manage Team",
            href: `/${organization.slug}/settings/teams/${teamIdInput}`,
          }}
          className="my-2"
        />
      ) : null}
      <div className="text-mauve-11 text-sm">
        <p>
          Use projects to isolate products that share nothing at all. Both data
          and setup is separate between projects.
        </p>
        <p className="italic">
          Required fields are marked with an asterisk (*).
        </p>
      </div>
      <div className="mt-4">
        <form className="flex gap-2" onSubmit={handleSubmit(onSubmit)}>
          <input
            className="border-mauve-6 placeholder:text-mauve-11 w-80 rounded-md border-1 px-2 py-1 text-sm"
            type="text"
            placeholder="Name*"
            {...register("name")}
          />
          <div className="relative w-52">
            <select
              {...register("teamId", { required: true })}
              name="teamId"
              className={cn(
                "border-mauve-6 h-full w-full appearance-none rounded-md border-1 px-2 py-1 text-sm",
                teamExists ? "" : "text-mauve-11",
              )}
              disabled={!teamExists}
            >
              {teamExists ? (
                <>
                  <option value="" disabled>
                    Team*
                  </option>
                  {teamsData.map((team, i) => (
                    <option key={i} value={team.id}>
                      {team.slug}
                    </option>
                  ))}
                </>
              ) : (
                <option value="" disabled>
                  Team*
                </option>
              )}
            </select>
            <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
              <ChevronDown width={15} className="stroke-mauve-10" />
            </div>
          </div>
          <div className="relative w-52">
            <select
              {...register("clusterId", { required: true })}
              name="clusterId"
              className={cn(
                "border-mauve-6 h-full w-full appearance-none rounded-md border-1 px-2 py-1 text-sm",
                clusterExists ? "" : "text-mauve-11",
              )}
              disabled={!clusterExists}
              defaultValue=""
            >
              {clusterExists ? (
                <>
                  <option value="" disabled>
                    Cluster*
                  </option>
                  {clustersData.map((cluster, i) => (
                    <option key={i} value={cluster.clusterId}>
                      {cluster.clusterName}
                    </option>
                  ))}
                </>
              ) : (
                <option value="" disabled>
                  Cluster*
                </option>
              )}
            </select>
            <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
              <ChevronDown width={15} className="stroke-mauve-10" />
            </div>
          </div>
          <Button
            className="w-32"
            disabled={!nameInput || !clusterIdInput || !teamIdInput}
            type="submit"
          >
            Create Project
          </Button>
        </form>
      </div>
    </div>
  );
}

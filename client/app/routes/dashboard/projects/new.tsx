import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import Button from "~/components/atoms/button/Button";
import { ChevronDown } from "~/components/atoms/icons";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { cn } from "~/utils/cn";
import { useTRPC } from "~/utils/trpc/react";

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

  useEffect(() => {
    setValue("teamId", teamsData?.[0] ? String(teamsData[0].id) : "");
  }, [teamsData?.[0]?.id, setValue]);

  useEffect(() => {
    setValue(
      "clusterId",
      clustersData?.[0] ? String(clustersData[0].clusterId) : "",
    );
  }, [clustersData?.[0]?.clusterId, setValue]);

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
      name: data.name,
      clusterId: Number(data.clusterId),
      teamId: Number(data.teamId),
    });
  };

  const clusterExists = clustersData && clustersData.length > 0;
  const teamExists = teamsData && teamsData.length > 0;
  const selectedTeamHasNoClusters =
    !!teamIdInput && teamExists && !clusterExists && !isLoading;

  return (
    <>
      <Breadcrumbs
        crumbs={[
          {
            label: (
              <button
                type="button"
                className="cursor-pointer hover:underline"
                onClick={() => navigate("../all", { relative: "path" })}
              >
                All Projects
              </button>
            ),
          },
          { label: "New Project" },
        ]}
      />
      <div className="flex flex-col gap-2 p-4">
        <h1 className="font-bold text-xl">New Project</h1>
        {!teamExists && !isLoading ? (
          <WarningBanner
            text="You must join or create a team before creating projects."
            linkOut={{
              text: "Manage Teams",
              href: `/${organization.slug}/settings/organization/teams`,
            }}
            className="my-2"
          />
        ) : null}
        {selectedTeamHasNoClusters ? (
          organization.isOwner ? (
            <WarningBanner
              text="This team has no clusters assigned. Create or assign a cluster before creating projects."
              linkOut={{
                text: "Create Cluster",
                href: `/${organization.slug}/clusters/new`,
              }}
              className="my-2"
            />
          ) : (
            <WarningBanner
              text="This team has no clusters assigned. Contact your admin to assign one before creating projects."
              className="my-2"
            />
          )
        ) : null}
        <div className="text-mauve-11 text-sm">
          <p>
            Use projects to isolate products that share nothing at all. Both
            data and setup is separate between projects.
          </p>
          <p className="italic">
            Required fields are marked with an asterisk (*).
          </p>
        </div>
        <div className="mt-4">
          <form className="flex gap-2" onSubmit={handleSubmit(onSubmit)}>
            <input
              className="w-80 rounded-md border border-mauve-6 px-2 py-1 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)] placeholder:text-mauve-11"
              type="text"
              placeholder="Name*"
              {...register("name")}
            />
            <div className="relative w-52">
              <select
                {...register("teamId", { required: true })}
                name="teamId"
                className={cn(
                  "h-full w-full appearance-none rounded-md border border-mauve-6 px-2 py-1 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]",
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
                  "h-full w-full appearance-none rounded-md border border-mauve-6 px-2 py-1 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]",
                  clusterExists ? "" : "text-mauve-11",
                )}
                disabled={!clusterExists}
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
    </>
  );
}

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
}

export default function NewProject() {
  const queryClient = useQueryClient();
  const trpc = useTRPC();
  const navigate = useNavigate();

  const organization = useOrganizationContext();

  const { data: clustersData, isLoading } = useQuery(
    trpc.organization.getOrganizationClusters.queryOptions({
      id: organization.id,
    }),
  );

  const createProjectMutation = useMutation(
    trpc.project.createProject.mutationOptions({
      onSuccess: async (project) => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${organization.slug}/projects/${project.id}`);
      },
    }),
  );

  const { register, handleSubmit, watch } = useForm<NewProjectFormInput>();
  const nameInput = watch("name", "");
  const clusterIdInput = watch("clusterId", "");

  const onSubmit: SubmitHandler<NewProjectFormInput> = (data) => {
    createProjectMutation.mutate({
      organizationId: organization.id,
      name: data.name,
      clusterId: Number(data.clusterId),
    });
  };

  const clusterExists = clustersData && clustersData.length > 0;

  return (
    <div className="flex flex-col gap-2 px-8 py-4">
      <h1 className="text-xl font-bold">New Project</h1>
      {!clusterExists && !isLoading ? (
        <WarningBanner
          text="You must create a cluster before creating projects."
          linkOut={{
            text: "Create Cluster",
            href: `/${organization.slug}/clusters/new`,
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
                clustersData.map((cluster, i) => (
                  <option key={i} value={cluster.id}>
                    {cluster.name}
                  </option>
                ))
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
            disabled={!nameInput || !clusterIdInput}
            type="submit"
          >
            Create Project
          </Button>
        </form>
      </div>
    </div>
  );
}

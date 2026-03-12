import Button from "~/components/atoms/button/Button";
import React from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useNavigate } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import { ChevronDown } from "~/components/atoms/icons";

interface NewClusterFormInput {
  name: string;
  serverType: string;
}

export default function NewCluster() {
  const queryClient = useQueryClient();
  const trpc = useTRPC();
  const navigate = useNavigate();

  const organization = useOrganizationContext();

  const { data: hetznerCredentialData, isLoading: isCredentialLoading } =
    useQuery(
      trpc.organization.getHetznerCredential.queryOptions({
        id: organization.id,
      }),
    );

  const isCredentialValid = !!hetznerCredentialData?.credential?.secret;

  const createClusterMutation = useMutation(
    trpc.cluster.createCluster.mutationOptions({
      onSuccess: async (newCluster) => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationClusters.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${organization.slug}/clusters/${newCluster.id}`);
      },
    }),
  );

  const { register, handleSubmit, watch } = useForm<NewClusterFormInput>();
  const nameInput = watch("name", "");

  const onSubmit: SubmitHandler<NewClusterFormInput> = (data) => {
    createClusterMutation.mutate({
      organizationId: organization.id,
      name: data.name,
      serverType: data.serverType,
    });
  };

  return (
    <div className="flex flex-col gap-2 px-8 py-4">
      <h1 className="text-xl font-bold">New Cluster</h1>
      {isCredentialLoading ? null : isCredentialValid ? null : (
        <WarningBanner
          text="You must enter your Hetzner API Key to create a cluster."
          linkOut={{
            text: "Organization Settings",
            href: `/${organization.slug}/settings/organization`,
          }}
          className="my-2"
        />
      )}
      <div className="text-mauve-11 text-sm">
        <p>
          A cluster is an isolated environment with its own compute resources,
          running independently.
        </p>
        <p className="italic">
          Required fields are marked with an asterisk (*).
        </p>
      </div>
      <div className="mt-4">
        <form className="flex gap-2" onSubmit={handleSubmit(onSubmit)}>
          <input
            className="border-mauve-6 w-80 rounded-md border-1 px-2 py-1 text-sm"
            type="text"
            placeholder="Name*"
            {...register("name")}
          />
          <div className="relative w-52">
            <select
              className="border-mauve-6 h-full w-full appearance-none rounded-md border-1 px-2 py-1 text-sm"
              defaultValue="cx23"
              {...register("serverType", { required: true })}
            >
              <option value="cx23">CX23 (2 CPU, 4 GB)</option>
              <option value="ccx33">CCX33 (8 CPU, 32 GB)</option>
            </select>
            <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
              <ChevronDown width={15} className="stroke-mauve-10" />
            </div>
          </div>
          <Button className="w-32" disabled={!nameInput} type="submit">
            Create Cluster
          </Button>
        </form>
      </div>
    </div>
  );
}

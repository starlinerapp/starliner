import React, { useEffect, useMemo } from "react";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useParams } from "react-router";
import Button from "~/components/atoms/button/Button";
import { useForm } from "react-hook-form";

interface FormInput {
  connectedBranch: string;
}

export default function UpdateConnectedBranchForm() {
  const trpc = useTRPC();
  const { id, environment } = useParams<{
    id: string;
    environment: string;
  }>();

  const { data: project } = useQuery(
    trpc.project.getProject.queryOptions({ id: Number(id) }),
  );

  const updateEnvironmentMutation = useMutation(
    trpc.environment.updateEnvironmentConnectedBranch.mutationOptions(),
  );

  const currentEnvironment = useMemo(
    () => project?.environments.find((e) => e.slug === environment),
    [project, environment],
  );

  const { data: connectedBranchData, isLoading: isConnectedBranchLoading } =
    useQuery(
      trpc.environment.getEnvironmentConnectedBranch.queryOptions({
        id: Number(currentEnvironment?.id),
      }),
    );

  const {
    register,
    handleSubmit,
    reset,
    formState: { isDirty },
  } = useForm<FormInput>({
    defaultValues: {
      connectedBranch: connectedBranchData?.branch,
    },
  });

  useEffect(() => {
    reset({ connectedBranch: connectedBranchData?.branch });
  }, [connectedBranchData?.branch, reset]);

  const onSubmit = (data: FormInput) => {
    updateEnvironmentMutation.mutate(
      {
        id: Number(currentEnvironment?.id),
        branchName: data.connectedBranch,
      },
      {
        onSuccess: async () => {
          reset({ connectedBranch: data.connectedBranch });
        },
      },
    );
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="flex items-center justify-between px-4 py-2">
        <div className="flex flex-col">
          <p className="text-md font-bold">Connected Branch</p>
          <p className="text-mauve-11 text-xs">
            Changes made to this GitHub branch will be automatically redeployed.
          </p>
        </div>
        {isConnectedBranchLoading ? (
          <Skeleton className="h-9.5 w-1/2" />
        ) : (
          <input
            className="border-mauve-6 disabled:text-mauve-11 w-1/2 rounded-md border-1 p-2"
            {...register("connectedBranch")}
          />
        )}
      </div>
      {isDirty && (
        <div className="flex justify-end gap-1 px-4 pb-2">
          <Button
            size="xs"
            className="w-20"
            intent="secondary"
            onClick={() => reset()}
          >
            Cancel
          </Button>
          <Button className="w-20" size="xs" type="submit">
            Save
          </Button>
        </div>
      )}
    </form>
  );
}

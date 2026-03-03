import Button from "~/components/atoms/button/Button";
import { ArrowRight } from "~/components/atoms/icons";
import React from "react";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { type SubmitHandler, useForm } from "react-hook-form";

interface DeployDatabaseFormInput {
  serviceName: string;
}

interface DeployDatabaseFormProps {
  defaultValues?: DeployDatabaseFormInput;
}

export default function DeployDatabaseForm({
  defaultValues,
}: DeployDatabaseFormProps) {
  const { register, handleSubmit, watch, reset } =
    useForm<DeployDatabaseFormInput>({
      defaultValues,
    });
  const serviceNameInput = watch("serviceName", "");

  const trpc = useTRPC();

  const createDatabaseMutation = useMutation(
    trpc.deployment.deployDatabase.mutationOptions(),
  );

  const { environment: currentEnvironment } = useEnvironment();

  const submit: SubmitHandler<DeployDatabaseFormInput> = async (data) => {
    if (!currentEnvironment) return;

    createDatabaseMutation.mutate(
      {
        id: currentEnvironment.id,
        serviceName: data.serviceName,
      },
      {
        onSuccess: () => {
          reset();
        },
      },
    );
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit(submit)}>
      <div className="flex flex-col gap-1">
        <p>PostgreSQL</p>
        <p className="text-mauve-11 truncate text-sm">
          Powerful, open source relational database
        </p>
      </div>
      <div className="flex flex-col gap-1">
        <p className="text-sm">Service Name</p>
        <div className="flex gap-2">
          <input
            className="border-mauve-6 disabled:text-mauve-10 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm disabled:hover:cursor-not-allowed"
            type="text"
            placeholder="Name*"
            disabled={!!defaultValues?.serviceName}
            {...register("serviceName", {
              required: true,
            })}
          />
        </div>
      </div>
      <Button
        type="submit"
        size="sm"
        className="w-28 flex-shrink-0 py-1.5"
        disabled={!!defaultValues || !serviceNameInput}
      >
        {defaultValues ? "Redeploy" : "Deploy"}
        <ArrowRight className="w-4 stroke-2" />
      </Button>
    </form>
  );
}

import Button from "~/components/atoms/button/Button";
import { ArrowRight } from "~/components/atoms/icons";
import React, { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

export interface DatabaseFormInput {
  serviceName: string;
}

interface DeployDatabaseFormProps {
  defaultValues?: DatabaseFormInput;
  onSubmit: (data: DatabaseFormInput) => Promise<void>;
  resetOnSuccess?: boolean;
}

export default function DeployDatabaseForm({
  defaultValues,
  onSubmit,
  resetOnSuccess = false,
}: DeployDatabaseFormProps) {
  const { register, handleSubmit, watch, reset } =
    useForm<DatabaseFormInput>({
      defaultValues,
    });
  const serviceNameInput = watch("serviceName", "");

  const [error, setError] = useState<string | null>(null);
  const isExistingDeployment = !!defaultValues;

  const submit: SubmitHandler<DatabaseFormInput> = async (data) => {
    try {
      await onSubmit(data);
      if (resetOnSuccess) {
        reset();
      }
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Oops something went wrong!");
    }
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit(submit)}>
      <div className="flex flex-col gap-1">
        <p>PostgreSQL</p>
        <p className="text-mauve-11 truncate text-sm">
          Powerful, open source relational database
        </p>
      </div>
      {error && (
        <div>
          <ErrorBanner text={error} />
        </div>
      )}
      <div className="flex flex-col gap-1">
        <p className="text-sm">Service Name</p>
        <div className="flex gap-2">
          <input
            className="border-mauve-6 disabled:text-mauve-10 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)] disabled:hover:cursor-not-allowed"
            type="text"
            placeholder="Name*"
            disabled={isExistingDeployment}
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
        disabled={!isExistingDeployment && !serviceNameInput}
      >
        {isExistingDeployment ? "Redeploy" : "Deploy"}
        <ArrowRight className="w-4 stroke-2" />
      </Button>
    </form>
  );
}

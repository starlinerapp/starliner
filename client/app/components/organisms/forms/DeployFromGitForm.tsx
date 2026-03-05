import React from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, Plus } from "~/components/atoms/icons";
import { type SubmitHandler, useFieldArray, useForm } from "react-hook-form";

export interface DeployFromGitFormInput {
  url: string;
  serviceName: string;
  dockerfilePath: string;
  projectDirectoryPath: string;
  port: number | null;
  envs: { name: string; value: string }[];
}

interface DeployFromGitFormProps {
  defaultValues?: DeployFromGitFormInput;
  onSubmit: (data: DeployFromGitFormInput) => Promise<void>;
  resetOnSuccess?: boolean;
}

export default function DeployFromGitForm({
  defaultValues,
  onSubmit,
  resetOnSuccess = false,
}: DeployFromGitFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    reset,
    control,
    formState: { isDirty },
  } = useForm<DeployFromGitFormInput>({ defaultValues });

  const { fields, append } = useFieldArray({
    control,
    name: "envs",
  });

  const urlInput = watch("url", "");
  const serviceNameInput = watch("serviceName", "");
  const port = watch("port", null);
  const projectDirectoryPathInput = watch("projectDirectoryPath", "");
  const dockerFilePathInput = watch("dockerfilePath", "");

  const submit: SubmitHandler<DeployFromGitFormInput> = async (data) => {
    data.envs = (data.envs ?? []).filter(
      (e) => e.name.trim() !== "" || e.value.trim() !== "",
    );

    await onSubmit(data);
    if (resetOnSuccess)
      reset({
        url: "",
        serviceName: "",
        dockerfilePath: "",
        projectDirectoryPath: "",
        port: null,
        envs: [],
      });
  };

  const inputValid =
    urlInput &&
    port &&
    serviceNameInput &&
    projectDirectoryPathInput &&
    dockerFilePathInput;

  return (
    <>
      <form className="flex flex-col gap-4" onSubmit={handleSubmit(submit)}>
        <div className="flex flex-col gap-1">
          <p>Import Git Repository</p>
          <p className="text-mauve-11 truncate text-sm">
            Enter a Git repository URL to deploy.
          </p>
        </div>
        <div className="flex flex-col gap-2">
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
          <div className="flex flex-col gap-1">
            <p className="text-sm">Git Repository URL</p>
            <div className="flex gap-2">
              <input
                className="border-mauve-6 disabled:text-mauve-10 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm disabled:hover:cursor-not-allowed"
                type="text"
                placeholder="Git URL*"
                disabled={!!defaultValues?.url}
                {...register("url", {
                  required: true,
                })}
              />
            </div>
          </div>
          <div className="flex flex-col gap-1">
            <div className="flex items-center gap-2">
              <div className="w-full">
                <p className="text-sm">Project Directory</p>
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
                  type="text"
                  placeholder="Path*"
                  {...register("projectDirectoryPath", {
                    required: true,
                  })}
                />
              </div>
              <div className="w-full">
                <p className="text-sm">Dockerfile</p>
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-24 rounded-md border-1 p-2 text-sm"
                  type="text"
                  placeholder="Path*"
                  {...register("dockerfilePath", {
                    required: true,
                  })}
                />
              </div>
            </div>
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-sm">Port</p>
            <div className="flex gap-2">
              <input
                className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
                type="number"
                placeholder="Port*"
                {...register("port", { required: true, valueAsNumber: true })}
              />
            </div>
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-sm">Environment Variables</p>
            {fields.map((field, index) => (
              <div key={field.id} className="flex gap-2">
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
                  type="text"
                  placeholder="Name*"
                  {...register(`envs.${index}.name`)}
                />
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
                  type="text"
                  placeholder="Value*"
                  {...register(`envs.${index}.value`)}
                />
              </div>
            ))}
            <Button
              intent="text"
              className="py-1"
              type="button"
              onClick={() => append({ name: "", value: "" })}
            >
              <Plus className="w-3 stroke-3" /> Add Another
            </Button>
          </div>
        </div>
        <Button
          type="submit"
          size="sm"
          disabled={!isDirty || !inputValid}
          className="w-28 flex-shrink-0 py-1.5"
        >
          {defaultValues ? "Redeploy" : "Deploy"}
          <ArrowRight className="w-4 stroke-2" />
        </Button>
      </form>
    </>
  );
}

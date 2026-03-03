import React from "react";
import { type SubmitHandler, useFieldArray, useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, Plus } from "~/components/atoms/icons";

interface EnvVar {
  name: string;
  value: string;
}

export interface ImageFormInput {
  serviceName: string;
  imageName: string;
  tag: string;
  port: number | null;
  envs: EnvVar[];
}

interface DeployImageFormProps {
  defaultValues?: ImageFormInput;
  onSubmit: (data: ImageFormInput) => Promise<void>;
  resetOnSuccess?: boolean;
}

export default function DeployImageForm({
  defaultValues,
  onSubmit,
  resetOnSuccess = false,
}: DeployImageFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    reset,
    control,
    formState: { isDirty },
  } = useForm<ImageFormInput>({
    defaultValues,
  });

  const { fields, append } = useFieldArray({
    control,
    name: "envs",
  });

  const serviceNameInput = watch("serviceName", "");
  const imageNameInput = watch("imageName", "");
  const tagInput = watch("tag", "");
  const portInput = watch("port", null);

  const submit: SubmitHandler<ImageFormInput> = async (data) => {
    data.envs = (data.envs ?? []).filter(
      (e) => e.name.trim() !== "" || e.value.trim() !== "",
    );

    await onSubmit(data);
    if (resetOnSuccess)
      reset({
        serviceName: "",
        imageName: "",
        tag: "",
        port: null,
        envs: [{ name: "", value: "" }],
      });
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit(submit)}>
      <div className="flex flex-col gap-1">
        <p>Docker Hub</p>
        <p className="text-mauve-11 truncate text-sm">
          Docker Hub Container Image Library
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
          <p className="text-sm">Image</p>
          <div className="flex items-center gap-2">
            <input
              className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
              type="text"
              placeholder="Name*"
              {...register("imageName", { required: true })}
            />
            {":"}
            <input
              className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-24 rounded-md border-1 p-2 text-sm"
              type="text"
              placeholder="Tag*"
              {...register("tag", { required: true })}
            />
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
        <Button
          size="sm"
          className="mt-2 w-28 flex-shrink-0 py-1.5"
          disabled={
            !isDirty ||
            !serviceNameInput ||
            !imageNameInput ||
            !tagInput ||
            !portInput
          }
        >
          {defaultValues ? "Redeploy" : "Deploy"}
          <ArrowRight className="w-4 stroke-2" />
        </Button>
      </div>
    </form>
  );
}

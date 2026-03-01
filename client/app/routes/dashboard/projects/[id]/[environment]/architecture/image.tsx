import React from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, Plus } from "~/components/atoms/icons";
import { type SubmitHandler, useFieldArray, useForm } from "react-hook-form";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";

interface EnvVar {
  name: string;
  value: string;
}

interface ImageFormInput {
  serviceName: string;
  imageName: string;
  tag: string;
  port: number | undefined;
  envs: EnvVar[];
}

export default function Image() {
  const trpc = useTRPC();

  const createImageMutation = useMutation(
    trpc.deployment.deployImage.mutationOptions(),
  );

  const { environment: currentEnvironment } = useEnvironment();

  const { register, handleSubmit, watch, reset, control } =
    useForm<ImageFormInput>();

  const { fields, append } = useFieldArray({
    control,
    name: "envs",
  });

  const serviceNameInput = watch("serviceName", "");
  const imageNameInput = watch("imageName", "");
  const tagInput = watch("tag", "");
  const portInput = watch("port", undefined);

  const onSubmit: SubmitHandler<ImageFormInput> = (data) => {
    if (!data.port) return;

    const sanitizedEnv = (data.envs ?? []).filter(
      (e) => e.name.trim() !== "" || e.value.trim() !== "",
    );

    createImageMutation.mutate(
      {
        id: currentEnvironment.id,
        serviceName: data.serviceName,
        imageName: data.imageName,
        tag: data.tag,
        port: data.port,
        envs: sanitizedEnv,
      },
      {
        onSuccess: () => {
          reset();
        },
      },
    );
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit(onSubmit)}>
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
              className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
              type="text"
              placeholder="Name*"
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
            !serviceNameInput || !imageNameInput || !tagInput || !portInput
          }
        >
          Deploy
          <ArrowRight className="w-4 stroke-2" />
        </Button>
      </div>
    </form>
  );
}

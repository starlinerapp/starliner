import React, { useState } from "react";
import { type SubmitHandler, useFieldArray, useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, Plus } from "~/components/atoms/icons";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { isEnvFile, parseEnvFile } from "~/service/envfile/envFile";

export interface ImageFormInput {
  serviceName: string;
  imageName: string;
  tag: string;
  port: number | null;
  volumeSizeMiB: number | null;
  volumeMountPath: string | null;
  envs: { name: string; value: string }[];
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
    formState: { errors },
  } = useForm<ImageFormInput>({
    defaultValues,
  });

  const { fields, append, replace } = useFieldArray({
    control,
    name: "envs",
  });

  const [error, setError] = useState<string | null>(null);
  const [addedVolume, setAddedVolume] = useState<boolean>(
    !!defaultValues?.volumeSizeMiB || !!defaultValues?.volumeMountPath,
  );

  const isExistingDeployment = !!defaultValues;
  const existingHasVolume =
    !!defaultValues?.volumeSizeMiB || !!defaultValues?.volumeMountPath;

  const serviceNameInput = watch("serviceName", "");
  const imageNameInput = watch("imageName", "");
  const tagInput = watch("tag", "");
  const portInput = watch("port", null);
  const volumeSizeInput = watch("volumeSizeMiB", null);
  const volumeMountPathInput = watch("volumeMountPath", null);

  const hasVolumeSize = !!volumeSizeInput && !isNaN(volumeSizeInput);
  const hasVolumeMountPath = !!volumeMountPathInput?.trim();
  const isVolumeIncomplete = hasVolumeSize !== hasVolumeMountPath;

  const submit: SubmitHandler<ImageFormInput> = async (data) => {
    data.envs = (data.envs ?? []).filter(
      (e) => e.name.trim() !== "" || e.value.trim() !== "",
    );

    if (!data.volumeSizeMiB || isNaN(data.volumeSizeMiB)) {
      data.volumeSizeMiB = null;
    }
    if (!data.volumeMountPath?.trim()) {
      data.volumeMountPath = null;
    }

    try {
      await onSubmit(data);
      if (resetOnSuccess)
        reset({
          serviceName: "",
          imageName: "",
          tag: "",
          port: null,
          volumeSizeMiB: null,
          volumeMountPath: null,
          envs: [],
        });
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Oops something went wrong!");
    }
  };

  const handleEnvPaste = (
    index: number,
    e: React.ClipboardEvent<HTMLInputElement>,
  ) => {
    const pasted = e.clipboardData.getData("text");
    if (!isEnvFile(pasted)) return;

    e.preventDefault();
    const parsed = parseEnvFile(pasted);
    if (parsed.length === 0) return;

    const before = fields
      .slice(0, index)
      .map((f) => ({ name: f.name, value: f.value }))
      .filter((f) => f.name.trim() !== "" || f.value.trim() !== "");

    replace([...before, ...parsed]);
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit(submit)}>
      <div className="flex flex-col gap-1">
        <p>Docker Hub</p>
        <p className="text-mauve-11 truncate text-sm">
          Docker Hub Container Image Library
        </p>
      </div>
      {error && (
        <div>
          <ErrorBanner text={error} />
        </div>
      )}
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
                onPaste={(e) => handleEnvPaste(index, e)}
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

        <div className="flex flex-col gap-1">
          <p className="text-sm">Persistent Volume</p>
          {isExistingDeployment && !existingHasVolume && (
            <p className="text-mauve-11 text-sm">No volume attached</p>
          )}
          {addedVolume && (
            <div className="flex flex-col gap-y-2">
              <div className="flex items-center gap-2">
                <input
                  className="border-mauve-6 disabled:text-mauve-10 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm disabled:hover:cursor-not-allowed"
                  type="number"
                  placeholder="Size*"
                  disabled={!!defaultValues?.volumeSizeMiB}
                  {...register("volumeSizeMiB", {
                    required: false,
                    valueAsNumber: true,
                  })}
                />
                <span className="text-mauve-12 text-sm">MiB</span>
              </div>

              <input
                className="border-mauve-6 disabled:text-mauve-10 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm disabled:hover:cursor-not-allowed"
                type="text"
                placeholder="Mount Path*"
                disabled={!!defaultValues?.volumeMountPath}
                {...register("volumeMountPath", {
                  required: false,
                  validate: (value) => {
                    if (!value) return true;
                    return (
                      /^\/(([a-zA-Z0-9_\-.]+)(\/[a-zA-Z0-9_\-.]+)*)?$/.test(
                        value,
                      ) || "Must be a valid absolute path (e.g. /data)"
                    );
                  },
                })}
              />
              {errors.volumeMountPath && (
                <p className="text-red-11 text-xs">
                  {errors.volumeMountPath.message}
                </p>
              )}
            </div>
          )}
          {!addedVolume && !isExistingDeployment && (
            <Button
              intent="text"
              className="py-1"
              type="button"
              onClick={() => setAddedVolume(true)}
            >
              <Plus className="w-3 stroke-3" /> Add Volume
            </Button>
          )}
        </div>
        <Button
          size="sm"
          className="mt-2 w-28 flex-shrink-0 py-1.5"
          disabled={
            !serviceNameInput ||
            !imageNameInput ||
            !tagInput ||
            !portInput ||
            isVolumeIncomplete
          }
        >
          {defaultValues ? "Redeploy" : "Deploy"}
          <ArrowRight className="w-4 stroke-2" />
        </Button>
      </div>
    </form>
  );
}

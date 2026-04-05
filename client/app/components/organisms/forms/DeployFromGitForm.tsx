import React, { useMemo, useState } from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, ChevronDown, Plus } from "~/components/atoms/icons";
import { type SubmitHandler, useFieldArray, useForm } from "react-hook-form";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import { isEnvFile, parseEnvFile } from "~/service/envfile/envFile";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import SelectProjectDirectoryDialog from "~/components/organisms/dialog/SelectProjectDirectoryDialog";
import SelectDockerfileDialog from "~/components/organisms/dialog/SelectDockerfileDialog";
import { cn } from "~/utils/cn";

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
  const trpc = useTRPC();
  const organization = useOrganizationContext();

  const { data: repositoriesData, isLoading } = useQuery(
    trpc.github.getRepositories.queryOptions({
      organizationId: organization.id,
    }),
  );

  const { register, handleSubmit, watch, reset, control, setValue } =
    useForm<DeployFromGitFormInput>({ defaultValues });

  const { fields, append, replace } = useFieldArray({
    control,
    name: "envs",
  });

  const [error, setError] = useState<string | null>(null);
  const [isSelectProjectDialogOpen, setIsSelectProjectDialogOpen] =
    useState(false);
  const [isSelectDockerFileDialogOpen, setIsSelectDockerFileDialogOpen] =
    useState(false);

  const { onChange: onUrlChange, ...urlRegister } = register("url", {
    required: true,
  });
  const urlInput = watch("url", "");
  const serviceNameInput = watch("serviceName", "");
  const port = watch("port", null);
  const projectDirectoryPathInput = watch("projectDirectoryPath", "");
  const dockerFilePathInput = watch("dockerfilePath", "");

  const selectedRepository = useMemo(() => {
    return repositoriesData?.find((repo) => repo.clone_url === urlInput);
  }, [repositoriesData, urlInput]);

  const submit: SubmitHandler<DeployFromGitFormInput> = async (data) => {
    data.envs = (data.envs ?? []).filter(
      (e) => e.name.trim() !== "" || e.value.trim() !== "",
    );

    try {
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
          <p>Select Git Repository</p>
          <p className="text-mauve-11 truncate text-sm">
            Select the repository you want to deploy.
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
            <p className="text-sm">Repository</p>
            <div className="flex gap-2">
              {isLoading ? (
                <Skeleton className="h-10 w-full" />
              ) : (
                <div className="relative w-full">
                  <select
                    {...urlRegister}
                    disabled={!!defaultValues?.url}
                    onChange={(e) => {
                      void onUrlChange(e);
                      setValue("projectDirectoryPath", "");
                      setValue("dockerfilePath", "");
                    }}
                    className={cn(
                      "border-mauve-6 bg-gray-2 hover:bg-gray-3 disabled:bg-gray-2 h-10 w-full cursor-pointer appearance-none rounded-md border-1 p-2 text-sm disabled:hover:cursor-not-allowed",
                      urlInput ? "text-mauve-12" : "text-mauve-11",
                    )}
                  >
                    <option value="">Select repository*</option>
                    {repositoriesData?.map((repo) => (
                      <option key={repo.clone_url} value={repo.clone_url}>
                        {repo.name}
                      </option>
                    ))}
                  </select>
                  <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                    <ChevronDown width={15} className="stroke-mauve-10" />
                  </div>
                </div>
              )}
            </div>
          </div>
          <div className="flex flex-col gap-1">
            <div className="flex items-center gap-2">
              <div className="w-full">
                <p className="text-sm">Project Directory</p>
                <div
                  className={cn(
                    "border-mauve-6 placeholder:text-mauve-11 bg-gray-2 hover:bg-gray-3 h-9.5 w-full min-w-52 cursor-pointer rounded-md border-1 p-2 text-sm",
                    !projectDirectoryPathInput && "text-mauve-11",
                    !selectedRepository && "hover:bg-gray-2 cursor-not-allowed",
                  )}
                  onClick={() => {
                    if (!selectedRepository) return;
                    setIsSelectProjectDialogOpen(true);
                  }}
                >
                  <p>{projectDirectoryPathInput || "Select directory*"}</p>
                </div>
                <input
                  type="hidden"
                  {...register("projectDirectoryPath", {
                    required: true,
                  })}
                />
              </div>
              <div className="w-full">
                <p className="text-sm">Dockerfile</p>
                <div
                  className={cn(
                    "border-mauve-6 placeholder:text-mauve-11 bg-gray-2 hover:bg-gray-3 h-9.5 w-full min-w-52 cursor-pointer rounded-md border-1 p-2 text-sm",
                    !dockerFilePathInput && "text-mauve-11",
                    (!selectedRepository || !projectDirectoryPathInput) &&
                      "hover:bg-gray-2 cursor-not-allowed",
                  )}
                  onClick={() => {
                    if (!selectedRepository || !projectDirectoryPathInput)
                      return;
                    setIsSelectDockerFileDialogOpen(true);
                  }}
                >
                  <p>{dockerFilePathInput || "Select Dockerfile*"}</p>
                </div>
                <input
                  type="hidden"
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
        </div>
        <Button
          type="submit"
          size="sm"
          disabled={!inputValid}
          className="w-28 flex-shrink-0 py-1.5"
        >
          {defaultValues ? "Redeploy" : "Deploy"}
          <ArrowRight className="w-4 stroke-2" />
        </Button>
      </form>
      <SelectProjectDirectoryDialog
        repositoryOwner={selectedRepository?.owner ?? ""}
        repositoryName={selectedRepository?.name ?? ""}
        isOpen={isSelectProjectDialogOpen}
        onOpenChange={setIsSelectProjectDialogOpen}
        onConfirm={(directory) => {
          setValue("projectDirectoryPath", directory);
          setValue("dockerfilePath", "");
        }}
      />
      <SelectDockerfileDialog
        repositoryOwner={selectedRepository?.owner ?? ""}
        repositoryName={selectedRepository?.name ?? ""}
        isOpen={isSelectDockerFileDialogOpen}
        onOpenChange={setIsSelectDockerFileDialogOpen}
        path={projectDirectoryPathInput}
        onConfirm={(file) => {
          setValue("dockerfilePath", file);
        }}
      />
    </>
  );
}

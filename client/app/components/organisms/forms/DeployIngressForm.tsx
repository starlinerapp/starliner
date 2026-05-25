import React, { useMemo, useState } from "react";
import { ArrowRight, ChevronDown, Minus, Plus } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import {
  type Control,
  type SubmitHandler,
  useFieldArray,
  useForm,
  type UseFormRegister,
} from "react-hook-form";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import {
  buildFullIngressHost,
  getIngressHostDomain,
  getIngressHostSuffix,
  isValidIngressHostPrefix,
  parseIngressHostPrefix,
} from "~/utils/ingress-host";

interface Path {
  path: string;
  pathType: "Prefix" | "Exact" | "";
  service: string;
}

interface Host {
  name: string;
  paths: Path[];
}

export interface IngressFormInput {
  hosts: Host[];
}

interface DeployIngressFormProps {
  deploymentEnvironment: string;
  defaultValues?: IngressFormInput;
  onSubmit: (data: IngressFormInput) => Promise<void>;
  resetOnSuccess?: boolean;
}

const emptyPathEntry: Path = { path: "", pathType: "", service: "" };
const emptyHostEntry: Host = { name: "", paths: [emptyPathEntry] };

export default function DeployIngressForm({
  deploymentEnvironment,
  defaultValues,
  onSubmit,
  resetOnSuccess = false,
}: DeployIngressFormProps) {
  const organization = useOrganizationContext();
  const { environment } = useEnvironment();

  const domain = getIngressHostDomain({
    deploymentEnvironment,
    environmentSlug: environment?.slug ?? "production",
  });

  const hostSuffix = getIngressHostSuffix(organization.slug, domain);

  const parsedDefaultValues = useMemo(() => {
    if (!defaultValues) {
      return undefined;
    }

    return {
      hosts: defaultValues.hosts.map((host) => ({
        ...host,
        name: parseIngressHostPrefix(host.name, organization.slug),
      })),
    };
  }, [defaultValues, organization.slug]);

  const {
    control,
    register,
    handleSubmit,
    watch,
    reset,
    formState: { isDirty, errors },
  } = useForm<IngressFormInput>({
    defaultValues: parsedDefaultValues ?? { hosts: [emptyHostEntry] },
    mode: "onSubmit",
  });

  const [error, setError] = useState<string | null>(null);

  const {
    fields: hostFields,
    append: appendHost,
    remove: removeHost,
  } = useFieldArray({
    control,
    name: "hosts",
  });

  const hosts = watch("hosts");

  const isFormIncomplete =
    !hosts ||
    hosts.length === 0 ||
    hosts.some((host) => {
      if (!isValidIngressHostPrefix(host.name)) return true;
      if (!host.paths || host.paths.length === 0) return true;

      return host.paths.some(
        (p) => !p.path?.trim() || !p.pathType?.trim() || !p.service?.trim(),
      );
    });

  const submit: SubmitHandler<IngressFormInput> = async (data) => {
    try {
      await onSubmit({
        hosts: data.hosts.map((host) => ({
          ...host,
          name: buildFullIngressHost(host.name, organization.slug, domain),
        })),
      });
      if (resetOnSuccess) {
        reset({
          hosts: [{ name: "", paths: [{ ...emptyPathEntry }] }],
        });
      }
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Oops something went wrong!");
    }
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit(submit)}>
      <div className="flex flex-col gap-1">
        <p>Traefik</p>
        <p className="text-mauve-11 truncate text-sm">
          Make your HTTP(S) network service available
        </p>
      </div>

      {error && (
        <div>
          <ErrorBanner text={error} />
        </div>
      )}

      <div className="flex flex-col gap-1">
        <p className="text-sm">Hosts</p>
        <div className="flex flex-col">
          {hostFields.map((host, hostIndex) => (
            <div key={host.id}>
              <HostEditor
                control={control}
                register={register}
                hostIndex={hostIndex}
                hostSuffix={hostSuffix}
                showRemove={hostIndex > 0}
                onRemove={() => removeHost(hostIndex)}
                errorMessage={errors.hosts?.[hostIndex]?.name?.message}
              />
            </div>
          ))}
        </div>

        <Button
          intent="text"
          type="button"
          onClick={() =>
            appendHost({
              name: "",
              paths: [{ ...emptyPathEntry }],
            })
          }
        >
          <Plus className="w-3 stroke-3" /> Add Host
        </Button>
      </div>

      <Button
        size="sm"
        className="w-28 shrink-0 py-1.5"
        type="submit"
        disabled={!isDirty || isFormIncomplete}
      >
        {defaultValues ? "Redeploy" : "Deploy"}
        <ArrowRight className="w-4 stroke-2" />
      </Button>
    </form>
  );
}

interface HostEditorProps {
  control: Control<IngressFormInput>;
  register: UseFormRegister<IngressFormInput>;
  hostIndex: number;
  hostSuffix: string;
  showRemove?: boolean;
  onRemove?: () => void;
  errorMessage?: string;
}

function HostEditor({
  control,
  register,
  hostIndex,
  hostSuffix,
  showRemove,
  onRemove,
  errorMessage,
}: HostEditorProps) {
  const trpc = useTRPC();

  const { environment } = useEnvironment();
  const { data: deploymentsData } = useQuery(
    trpc.environment.getEnvironmentDeployments.queryOptions({
      id: environment?.id,
    }),
  );

  const {
    fields: pathFields,
    append: appendPath,
    remove: removePath,
  } = useFieldArray({
    control,
    name: `hosts.${hostIndex}.paths`,
  });

  const linkableServices = [
    ...(deploymentsData?.gitDeployments.map((d) => d.serviceName) ?? []),
    ...(deploymentsData?.images.map((d) => d.serviceName) ?? []),
  ];

  return (
    <div className="relative">
      {showRemove && (
        <button
          type="button"
          onClick={onRemove}
          className="border-mauve-6 hover:border-violet-9 hover:text-violet-9 text-mauve-7 absolute top-3.75 -left-2 z-10 flex h-4 w-4 cursor-pointer items-center justify-center rounded border-[1.5px] bg-white"
        >
          <Minus className="stroke-3" />
        </button>
      )}

      <div className="border-mauve-6 relative flex flex-col gap-1 border-l-2 pl-6">
        <div className="flex gap-2">
          <div className="border-mauve-6 absolute -left-0.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />
          <div className="flex w-full min-w-0 flex-col gap-1">
            <div className="flex w-full min-w-0 items-center gap-1">
              <input
                className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 min-w-24 flex-1 rounded-md border p-2 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                type="text"
                placeholder="Subdomain*"
                {...register(`hosts.${hostIndex}.name`, {
                  validate: (value) =>
                    isValidIngressHostPrefix(value) ||
                    `Use lowercase letters, numbers, and hyphens before${hostSuffix}`,
                })}
              />
              <span className="text-mauve-11 truncate text-sm">
                {hostSuffix}
              </span>
            </div>
            {errorMessage && (
              <p className="text-red-11 text-xs">{errorMessage}</p>
            )}
          </div>
        </div>

        <div className="border-mauve-6 relative flex flex-col gap-3 border-l-2 pl-6">
          {pathFields.map((path, pathIndex) => (
            <div key={path.id} className="relative flex flex-col gap-1">
              {pathIndex > 0 && (
                <button
                  type="button"
                  onClick={() => removePath(pathIndex)}
                  className="border-mauve-6 hover:border-violet-9 hover:text-violet-9 text-mauve-7 absolute top-3.75 -left-8 z-10 flex h-4 w-4 cursor-pointer items-center justify-center rounded border-[1.5px] bg-white"
                >
                  <Minus className="stroke-3" />
                </button>
              )}

              <div className="border-mauve-6 absolute -left-6.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />

              <div className="flex w-full gap-1">
                <div className="relative min-w-24">
                  <select
                    className="border-mauve-6 bg-gray-2 h-full w-full appearance-none rounded-md border px-2 py-1 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                    defaultValue={path.pathType ?? ""}
                    {...register(
                      `hosts.${hostIndex}.paths.${pathIndex}.pathType`,
                    )}
                  >
                    <option value="" disabled>
                      Type*
                    </option>
                    {["Prefix", "Exact"].map((pathType) => (
                      <option key={pathType} value={pathType}>
                        {pathType}
                      </option>
                    ))}
                  </select>

                  <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                    <ChevronDown width={15} className="stroke-mauve-10" />
                  </div>
                </div>

                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-24 rounded-md border p-2 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                  type="text"
                  placeholder="Path*"
                  {...register(`hosts.${hostIndex}.paths.${pathIndex}.path`)}
                />
              </div>

              <div className="border-mauve-6 relative ml-1 flex flex-col gap-1 border-l-2 pl-6">
                <div className="border-mauve-6 absolute -left-0.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />

                <div className="relative min-w-48">
                  <select
                    className="border-mauve-6 bg-gray-2 h-full w-full appearance-none rounded-md border p-2 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                    defaultValue={path.service ?? ""}
                    {...register(
                      `hosts.${hostIndex}.paths.${pathIndex}.service`,
                    )}
                  >
                    <option value="" disabled>
                      Service*
                    </option>
                    {linkableServices.map((name, i) => (
                      <option key={i} value={name}>
                        {name}
                      </option>
                    ))}
                  </select>

                  <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                    <ChevronDown width={15} className="stroke-mauve-10" />
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        <Button
          intent="text"
          className="py-1"
          type="button"
          onClick={() => appendPath({ ...emptyPathEntry })}
        >
          <Plus className="w-3 stroke-3" /> Add Path
        </Button>
      </div>
    </div>
  );
}

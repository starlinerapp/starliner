import React from "react";
import { ArrowRight, ChevronDown, Minus, Plus } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";
import { cn } from "~/utils/cn";
import {
  useFieldArray,
  useForm,
  type Control,
  type UseFormRegister,
} from "react-hook-form";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

interface Path {
  path: string;
  pathType: "Prefix" | "Exact" | "";
  service: string;
}

interface Host {
  name: string;
  paths: Path[];
}

interface IngressFormInput {
  hosts: Host[];
}

const emptyPathEntry: Path = { path: "", pathType: "", service: "" };
const emptyHostEntry: Host = { name: "", paths: [emptyPathEntry] };

export default function Ingress() {
  const trpc = useTRPC();
  const createIngressMutation = useMutation(
    trpc.deployment.deployIngress.mutationOptions(),
  );
  const { environment: currentEnvironment, clusterId } = useEnvironment();
  const { data: clusterData } = useQuery(
    trpc.cluster.getCluster.queryOptions(
      { id: clusterId! },
      { enabled: !!clusterId },
    ),
  );

  const { control, register, handleSubmit, watch, reset } =
    useForm<IngressFormInput>({
      defaultValues: { hosts: [emptyHostEntry] },
      mode: "onSubmit",
    });

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
      if (!host.name?.trim()) return true;
      if (!host.paths || host.paths.length === 0) return true;

      return host.paths.some(
        (p) => !p.path?.trim() || !p.pathType?.trim() || !p.service?.trim(),
      );
    });

  const onSubmit = (data: IngressFormInput) => {
    if (!currentEnvironment) return;
    createIngressMutation.mutate(
      {
        id: currentEnvironment.id,
        ingressHosts: data.hosts.map((h) => ({
          host: h.name + `.${clusterData?.ipv4Address}.nip.io`,
          paths: h.paths.map((p) => ({
            path: p.path,
            pathType: p.pathType as "Prefix" | "Exact",
            serviceName: p.service,
          })),
        })),
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
        <p>Traefik</p>
        <p className="text-mauve-11 truncate text-sm">
          Make your HTTP(S) network service available
        </p>
      </div>

      <div className="flex flex-col gap-1">
        <p className="text-sm">Hosts</p>

        <div className="flex flex-col">
          {hostFields.map((host, hostIndex) => (
            <div key={host.id}>
              <HostEditor
                control={control}
                register={register}
                hostIndex={hostIndex}
                showRemove={hostIndex > 0}
                onRemove={() => removeHost(hostIndex)}
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
        className="w-28 flex-shrink-0 py-1.5"
        type="submit"
        disabled={isFormIncomplete || createIngressMutation.isPending}
      >
        Deploy
        <ArrowRight className="w-4 stroke-2" />
      </Button>
    </form>
  );
}

interface HostEditorProps {
  control: Control<IngressFormInput>;
  register: UseFormRegister<IngressFormInput>;
  hostIndex: number;
  showRemove?: boolean;
  onRemove?: () => void;
}

function HostEditor({
  control,
  register,
  hostIndex,
  showRemove,
  onRemove,
}: HostEditorProps) {
  const trpc = useTRPC();

  const { environment, clusterId } = useEnvironment();
  const { data: clusterData } = useQuery(
    trpc.cluster.getCluster.queryOptions(
      { id: clusterId! },
      { enabled: !!clusterId },
    ),
  );

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

  return (
    <div className="relative">
      {showRemove && (
        <button
          type="button"
          onClick={onRemove}
          className="border-mauve-6 hover:border-violet-9 hover:text-violet-9 text-mauve-7 absolute top-[15px] -left-2 z-10 flex h-4 w-4 cursor-pointer items-center justify-center rounded border-[1.5px] bg-white"
        >
          <Minus className="stroke-3" />
        </button>
      )}

      <div className="border-mauve-6 relative flex flex-col gap-1 border-l-2 pl-6">
        <div className="flex gap-2">
          <div className="border-mauve-6 absolute -left-0.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />
          <div className="flex w-full items-center gap-1">
            <input
              className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
              type="text"
              placeholder="Host*"
              {...register(`hosts.${hostIndex}.name`)}
            />
            {clusterData?.ipv4Address ? (
              <div className="border-mauve-6 text-mauve-11 flex cursor-not-allowed rounded-md border-1 p-2 text-sm">
                <p>.{clusterData?.ipv4Address}</p>
                <p>.nip.io</p>
              </div>
            ) : (
              <Skeleton className="h-full w-56" />
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
                  className="border-mauve-6 hover:border-violet-9 hover:text-violet-9 text-mauve-7 absolute top-[15px] -left-8 z-10 flex h-4 w-4 cursor-pointer items-center justify-center rounded border-[1.5px] bg-white"
                >
                  <Minus className="stroke-3" />
                </button>
              )}

              <div className="border-mauve-6 absolute -left-6.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />

              <div className="flex w-full gap-1">
                <div className="relative min-w-24">
                  <select
                    className={cn(
                      "border-mauve-6 bg-gray-2 h-full w-full appearance-none rounded-md border-1 px-2 py-1 text-sm",
                      "text-mauve-11",
                    )}
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
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-24 rounded-md border p-2 text-sm"
                  type="text"
                  placeholder="Path*"
                  {...register(`hosts.${hostIndex}.paths.${pathIndex}.path`)}
                />
              </div>

              <div className="border-mauve-6 relative ml-1 flex flex-col gap-1 border-l-2 pl-6">
                <div className="border-mauve-6 absolute -left-0.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />

                <div className="relative min-w-48">
                  <select
                    className={cn(
                      "border-mauve-6 bg-gray-2 h-full w-full appearance-none rounded-md border-1 p-2 text-sm",
                      "text-mauve-11",
                    )}
                    defaultValue={path.service ?? ""}
                    {...register(
                      `hosts.${hostIndex}.paths.${pathIndex}.service`,
                    )}
                  >
                    <option value="" disabled>
                      Service*
                    </option>
                    {deploymentsData?.images.map((d, i) => (
                      <option key={i} value={d.serviceName}>
                        {d.serviceName}
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

import React, { useMemo } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/atoms/dropdown/DropdownMenu";
import ChevronDown from "~/components/atoms/icons/ChevronDown";
import { ChevronRight } from "~/components/atoms/icons";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import { useNavigate, useParams } from "react-router";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { type SubmitHandler, useForm } from "react-hook-form";
import type { ResponseProject } from "~/server/api/client/generated";
import { cn } from "~/utils/cn";

interface NewEnvironmentFormInput {
  environmentName: string;
  creationMode: "duplicate" | "empty";
  sourceEnvironment?: number;
}

interface ManageEnvironmentsProps {
  organization: {
    id: number;
    name: string;
    slug: string;
    isOwner: boolean;
  };
  project: ResponseProject | undefined;
}

export default function ManageEnvironments({
  project,
  organization,
}: ManageEnvironmentsProps) {
  const navigate = useNavigate();
  const { slug, id, environment } = useParams<{
    slug: string;
    id: string;
    environment: string;
  }>();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const createEnvironmentMutation = useMutation(
    trpc.environment.createEnvironment.mutationOptions({
      onSuccess: async (data) => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getUserProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${slug}/projects/${id}/${data.slug}/architecture/git`);
      },
    }),
  );

  const environments = project?.environments.map((env) => env) ?? [];
  const currentEnvironment = useMemo(
    () => project?.environments.find((e) => e.slug === environment),
    [project, environment],
  );

  const [environmentDialogOpen, setEnvironmentDialogOpen] =
    React.useState(false);
  const {
    register,
    handleSubmit,
    watch,
    formState: { isValid },
  } = useForm<NewEnvironmentFormInput>({
    defaultValues: {
      environmentName: "",
      creationMode: "duplicate",
    },
  });
  const creationMode = watch("creationMode");

  const onNewEnvironmentSubmit: SubmitHandler<NewEnvironmentFormInput> = (
    data,
  ) => {
    createEnvironmentMutation.mutate({
      name: data.environmentName,
      organizationId: organization.id,
      projectId: project?.id ?? 0,
      sourceEnvironmentId:
        creationMode === "duplicate" ? data.sourceEnvironment : undefined,
    });
    setEnvironmentDialogOpen(false);
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <div className="border-violet-10 flex cursor-pointer items-center gap-1.5 rounded-md border-[1px] px-2 text-sm">
            <h1>{currentEnvironment?.name}</h1>
            <ChevronDown height={15} width={15} />
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="min-w-48">
          {environments.length > 0 &&
            environments
              .filter((env) => currentEnvironment?.slug !== env.slug)
              .map((env) => (
                <DropdownMenuItem
                  className="text-xs"
                  key={env.slug}
                  onClick={() => {
                    navigate(
                      `/${slug}/projects/${id}/${env.slug}/architecture/git`,
                    );
                  }}
                >
                  <div>{env.name}</div>
                </DropdownMenuItem>
              ))}
          {environments.length > 1 && <DropdownMenuSeparator />}
          <DropdownMenuItem onClick={() => setEnvironmentDialogOpen(true)}>
            <div className="flex w-full items-center justify-between gap-1.5">
              <p className="text-xs">New Environment</p>
              <ChevronRight width={12} height={12} />
            </div>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <Dialog
        onOpenChange={(value) => setEnvironmentDialogOpen(value)}
        open={environmentDialogOpen}
      >
        <DialogContent>
          <form
            className="flex flex-col gap-5"
            onSubmit={handleSubmit(onNewEnvironmentSubmit)}
          >
            <div className="flex flex-col gap-4">
              <div className="flex flex-col gap-2">
                <h1>New Environment</h1>
                <p className="text-mauve-11 text-sm">
                  All the changes will be isolated from other environments.
                </p>
              </div>
              <input
                type="text"
                placeholder="Name*"
                className="border-mauve-6 rounded-md border-1 p-2 text-sm focus:outline-none"
                {...register("environmentName", {
                  required: true,
                  validate: (v) => v.trim().length > 0,
                })}
              />
              <div className="flex flex-col gap-2">
                <label
                  htmlFor="creation-mode-duplicate"
                  className={cn(
                    "border-mauve-6 flex cursor-pointer items-start gap-2 rounded-md border p-2",
                    creationMode === "duplicate" &&
                      "border-violet-9 bg-violet-2",
                  )}
                >
                  <input
                    id="creation-mode-duplicate"
                    type="radio"
                    value="duplicate"
                    className="mt-[3px]"
                    {...register("creationMode")}
                  />
                  <div className="flex w-full flex-col gap-2">
                    <h2 className="text-sm">Duplicate Environment</h2>
                    <div className="relative w-full">
                      <select
                        {...register("sourceEnvironment", {
                          required: creationMode === "duplicate",
                          valueAsNumber: true,
                        })}
                        className="border-mauve-6 hover:bg-gray-3 h-10 w-full cursor-pointer appearance-none rounded-md border bg-white p-2 text-sm"
                        onClick={(e) => e.stopPropagation()}
                      >
                        {environments.map((env) => (
                          <option key={env.id} value={env.id}>
                            {env.name}
                          </option>
                        ))}
                      </select>
                      <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                        <ChevronDown width={15} className="stroke-mauve-10" />
                      </div>
                    </div>
                    <p className="text-mauve-11 text-sm">
                      Copy all the services from an existing environment.
                    </p>
                  </div>
                </label>

                <label
                  htmlFor="creation-mode-empty"
                  className={cn(
                    "border-mauve-6 flex cursor-pointer items-start gap-2 rounded-md border p-2",
                    creationMode === "empty" && "border-violet-9 bg-violet-2",
                  )}
                >
                  <input
                    id="creation-mode-empty"
                    type="radio"
                    value="empty"
                    className="mt-[3px]"
                    {...register("creationMode")}
                  />
                  <div className="flex w-full flex-col gap-2">
                    <h2 className="text-sm">Empty Environment</h2>
                    <p className="text-mauve-11 text-sm">
                      An empty environment with no services included.
                    </p>
                  </div>
                </label>
              </div>
            </div>
            <div className="flex justify-end gap-2">
              <Button
                intent="secondary"
                className="w-24"
                onClick={() => setEnvironmentDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                className="w-24 self-end"
                disabled={!isValid}
              >
                Create
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}

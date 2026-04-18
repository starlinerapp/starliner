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

interface NewEnvironmentFormInput {
  environmentName: string;
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
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getUserProjects.queryKey({
            id: organization.id,
          }),
        });
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
    formState: { isDirty },
  } = useForm<NewEnvironmentFormInput>({
    defaultValues: {
      environmentName: "",
    },
  });

  const onNewEnvironmentSubmit: SubmitHandler<NewEnvironmentFormInput> = (
    data,
  ) => {
    createEnvironmentMutation.mutate({
      name: data.environmentName,
      organizationId: organization.id,
      projectId: project?.id ?? 0,
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
        <DropdownMenuContent>
          {environments.length > 0 &&
            environments
              .filter((env) => currentEnvironment?.slug !== env.slug)
              .map((env) => (
                <DropdownMenuItem
                  className="text-xs"
                  key={env.slug}
                  onClick={() => {
                    navigate(
                      `/${slug}/projects/${id}/${env.slug}/architecture`,
                    );
                  }}
                >
                  <div>{env.name}</div>
                </DropdownMenuItem>
              ))}
          {environments.length > 1 && <DropdownMenuSeparator />}
          <DropdownMenuItem onClick={() => setEnvironmentDialogOpen(true)}>
            <div className="flex items-center gap-1.5">
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
                placeholder="Staging"
                className="border-mauve-6 rounded-md border-1 px-2 py-1 text-sm focus:outline-none"
                {...register("environmentName")}
              />
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
                disabled={!isDirty}
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

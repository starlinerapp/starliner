import LinkNavigationBar from "~/components/organisms/navigation-bar/LinkNavigationBar";
import { Outlet, useNavigate, useParams } from "react-router";
import React, { useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/atoms/dropdown/DropdownMenu";
import { ChevronDown, ChevronRight } from "~/components/atoms/icons";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import { type SubmitHandler, useForm } from "react-hook-form";

interface NewEnvironmentFormInput {
  environmentName: string;
}

export default function ProjectLayout() {
  const { slug, id, environment } = useParams<{
    slug: string;
    id: string;
    environment: string;
  }>();
  const organization = useOrganizationContext();

  const queryClient = useQueryClient();
  const trpc = useTRPC();

  const navigate = useNavigate();

  const { data: projects, isLoading } = useQuery(
    trpc.organization.getOrganizationProjects.queryOptions({
      id: organization.id,
    }),
  );
  const createEnvironmentMutation = useMutation(
    trpc.environment.createEnvironment.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationProjects.queryKey({
            id: organization.id,
          }),
        });
      },
    }),
  );

  const currentProject = projects?.find((p) => p.id === Number(id));

  const navigationBarItems = [
    {
      title: "Architecture",
      href: `/${slug}/projects/${id}/${environment}/architecture`,
    },
    {
      title: "Observability",
      href: `/${slug}/projects/${id}/${environment}/observability`,
    },
    { title: "Logs", href: `/${slug}/projects/${id}/${environment}/logs` },
    {
      title: "Settings",
      href: `/${slug}/projects/${id}/${environment}/settings`,
    },
  ];

  const environments = currentProject?.environments.map((env) => env) ?? [];

  const [selectedEnvironment, setSelectedEnvironment] = React.useState<
    string | undefined
  >(environment);

  const [environmentDialogOpen, setEnvironmentDialogOpen] =
    React.useState(false);

  useEffect(() => {
    if (environments.length > 0 && !selectedEnvironment)
      setSelectedEnvironment(environments[0].slug);
  }, [environments]);

  const { register, handleSubmit } = useForm<NewEnvironmentFormInput>();

  const onNewEnvironmentSubmit: SubmitHandler<NewEnvironmentFormInput> = (
    data,
  ) => {
    createEnvironmentMutation.mutate({
      name: data.environmentName,
      organizationId: organization.id,
      projectId: currentProject?.id ?? 0,
    });
    setEnvironmentDialogOpen(false);
  };

  return (
    <div className="bg-violet-1 flex h-full flex-col">
      {isLoading ? (
        <div className="bg-violet-1 px-4 pt-4">
          <Skeleton className="h-7 w-32" />
        </div>
      ) : (
        <div className="bg-violet-1 flex items-center gap-3 px-4 pt-4">
          <h1 className="text-mauve-12 text-xl font-bold">
            {currentProject?.name}
          </h1>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <div className="border-violet-10 flex cursor-pointer items-center gap-1.5 rounded-md border-[1px] px-2 text-sm">
                <h1>
                  {
                    environments.find((e) => e.slug === selectedEnvironment)
                      ?.name
                  }
                </h1>
                <ChevronDown height={15} width={15} />
              </div>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              {environments.length > 0 &&
                environments
                  .filter((env) => selectedEnvironment !== env.slug)
                  .map((env) => (
                    <DropdownMenuItem
                      className="text-xs"
                      key={env.slug}
                      onClick={() => {
                        setSelectedEnvironment(env.slug);
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
                    <p className="text-mauve-11 text-xs">
                      All the changes will be isolated from other environments,
                      you can sync environments to pass changes.
                    </p>
                  </div>
                  <input
                    type="text"
                    placeholder="Staging"
                    className="border-mauve-6 rounded-md border-1 px-2 py-1 text-sm focus:outline-none"
                    {...register("environmentName")}
                  />
                </div>
                <Button type="submit" className="w-40 self-end">
                  Create Environment
                </Button>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      )}
      <LinkNavigationBar items={navigationBarItems} />
      <Outlet />
    </div>
  );
}

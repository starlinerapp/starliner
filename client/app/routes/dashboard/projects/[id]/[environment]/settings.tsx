import React, {useEffect} from "react";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import {type SubmitHandler, useForm} from "react-hook-form";
import Skeleton from "~/components/atoms/skeleton/Skeleton";


interface ProjectSettingsForm {
  projectName: string;
}

export default function ProjectSettings() {
  const navigate = useNavigate();
  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const { slug, id } = useParams<{
    slug: string;
    id: string;
  }>();

  const {data: project, isLoading, error} = useQuery(trpc.project.getProject.queryOptions({
    id: Number(id)
  }))

  const { register, handleSubmit, reset, watch } = useForm<ProjectSettingsForm>({});

  const projectNameInput = watch("projectName", "")

  useEffect(() => {
    if (project) {
      reset({
        projectName: project.name,
      });
    }
  }, [project, reset]);


  useEffect( () => {
    if(!project) {
      return;
    }

    const timeoutId = setTimeout(() => {
      if (projectNameInput.trim()) {
        onSubmit({ projectName: projectNameInput });
      }
    }, 1000);

    return () => clearTimeout(timeoutId);
  }, [projectNameInput]);

  const onSubmit: SubmitHandler<ProjectSettingsForm> = (data) => {
    updateProjectMutation.mutate({projectId: Number(id), name: data.projectName})
  }

  const updateProjectMutation = useMutation(
      trpc.project.updateProjectName.mutationOptions({
        onSuccess: async () => {
          await queryClient.invalidateQueries({
            queryKey: trpc.project.getProject.queryKey({
              id: Number(id)
            })
          })
          await queryClient.invalidateQueries({
            queryKey: trpc.organization.getOrganizationProjects.queryKey({
              id: organization.id,
            }),
          })
        }
      })
  )

  const deleteProjectMutation = useMutation(
    trpc.project.deleteProject.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${slug}/projects/all`);
      },
    }),
  );

  return (
      <div className="flex flex-col gap-4 w-full p-4 xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs uppercase">
            General
          </div>
          <form className="flex items-center justify-between px-4 py-2" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <p className="text-md font-bold">Project Name</p>
              <p className="text-mauve-11 text-xs">
                A human friendly name for the project.
              </p>
            </div>
            {isLoading ? (
             <Skeleton className="h-7 w-80" />
            ) : (
              <input
                  className="border-mauve-6 placeholder:text-mauve-11 w-80 rounded-md border-1 px-2 py-1 text-sm"
                  type="text"
                  placeholder="Name*"
                  {...register("projectName")}
              />
            )}
          </form>
        </div>
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs uppercase">
            Danger Zone
          </div>
          <div className="flex items-center justify-between px-4 py-2">
            <div>
              <p className="text-md font-bold">Delete this Project</p>
              <p className="text-mauve-11 text-xs">
                Once you delete a project, there is no going back. Please be
                certain.
              </p>
            </div>
            <Button
                className="w-36"
                intent="danger"
                size="sm"
                onClick={() =>
                    deleteProjectMutation.mutate({
                      id: Number(id),
                    })
                }
            >
              Delete this Project
            </Button>
          </div>
        </div>
      </div>
  );
}

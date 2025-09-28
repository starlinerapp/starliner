import React from "react";
import Button from "~/components/atoms/button/Button";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useNavigate } from "react-router";

interface NewProjectFormInput {
  name: string;
}

export default function NewProject() {
  const queryClient = useQueryClient();
  const trpc = useTRPC();
  const navigate = useNavigate();

  const organization = useOrganizationContext();

  const createProjectMutation = useMutation(
    trpc.project.createProject.mutationOptions({
      onSuccess: async () => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationProjects.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${organization.slug}/projects/all`);
      },
    }),
  );

  const { register, handleSubmit, watch } = useForm<NewProjectFormInput>();
  const nameInput = watch("name", "");

  const onSubmit: SubmitHandler<NewProjectFormInput> = (data) => {
    createProjectMutation.mutate({
      organizationId: organization.id,
      name: data.name,
    });
  };

  return (
    <div className="flex flex-col gap-2 px-8 py-4">
      <h1 className="text-xl font-bold">New Project</h1>
      <div className="text-mauve-11 text-sm">
        <p>
          Use projects to isolate products that share nothing at all. Both data
          and setup is separate between projects.
        </p>
        <p className="italic">
          Required fields are marked with an asterisk (*).
        </p>
      </div>
      <div className="mt-4">
        <form className="flex gap-2" onSubmit={handleSubmit(onSubmit)}>
          <input
            className="border-mauve-6 w-80 rounded-md border-1 px-2 py-1 text-sm"
            type="text"
            placeholder="Name*"
            {...register("name")}
          />
          <Button className="w-32" disabled={!nameInput} type="submit">
            Create Project
          </Button>
        </form>
      </div>
    </div>
  );
}

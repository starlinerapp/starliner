import { useMutation } from "@tanstack/react-query";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import { getAuthClient } from "~/utils/auth/client";
import { useTRPC } from "~/utils/trpc/react";

interface NewOrganizationFormInput {
  name: string;
}

export default function NewOrganization() {
  const authClient = getAuthClient();
  const trpc = useTRPC();
  const navigate = useNavigate();
  const { register, handleSubmit } = useForm<NewOrganizationFormInput>();
  const createOrganizationMutation = useMutation(
    trpc.organization.createOrganization.mutationOptions(),
  );

  const onSubmit: SubmitHandler<NewOrganizationFormInput> = async (data) => {
    createOrganizationMutation.mutate(
      { name: data.name },
      {
        onSuccess: (org) => navigate(`/organizations/${org.slug}/githubapp`),
      },
    );
  };

  async function handleSignOutClicked() {
    await authClient.signOut();
    navigate("/login");
  }

  return (
    <div className="flex w-125 flex-col gap-4">
      <button
        type="button"
        className="flex cursor-pointer items-center gap-0.5 self-end rounded-md px-2 py-0.5 font-light text-sm hover:bg-gray-4"
        onClick={handleSignOutClicked}
      >
        Sign out <ArrowRight className="w-3" />
      </button>
      <h1 className="font-medium text-xl">Create a New Organization</h1>
      <p className="text-mauve-11 text-sm">
        Organizations represent the top level in your hierarchy. You&#39;ll be
        able to bundle a collection of teams within an organization as well as
        give organization-wide permissions to users.
      </p>
      <form className="flex flex-col gap-2" onSubmit={handleSubmit(onSubmit)}>
        <span className="flex flex-col gap-1">
          <label htmlFor="name" className="text-sm">
            Organization Name
          </label>
          <input
            className="rounded-md border-1 border-mauve-6 p-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
            type="text"
            placeholder="e.g. My Company"
            {...register("name")}
          />
        </span>
        <Button type="submit" className="mt-2" size="md">
          Create Organization <ChevronRight className="w-4 stroke-3" />
        </Button>
      </form>
    </div>
  );
}

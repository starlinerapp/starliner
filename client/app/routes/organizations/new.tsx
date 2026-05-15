import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import React from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { useAuthClient } from "~/utils/auth/client";
import { useMutation } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";

interface NewOrganizationFormInput {
  name: string;
}

export default function NewOrganization() {
  const authClient = useAuthClient();
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
        onSuccess: (org) => navigate(`/${org.slug}`),
      },
    );
  };

  async function handleSignOutClicked() {
    await authClient.signOut();
    navigate("/login");
  }

  return (
    <div className="flex w-[500px] flex-col gap-4">
      <button
        className="hover:bg-gray-4 flex cursor-pointer items-center gap-0.5 self-end rounded-md px-2 py-0.5 text-sm font-light"
        onClick={handleSignOutClicked}
      >
        Sign out <ArrowRight className="w-3" />
      </button>
      <h1 className="text-xl font-medium">Create a New Organization</h1>
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
            className="border-mauve-6 rounded-md border-1 p-2"
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

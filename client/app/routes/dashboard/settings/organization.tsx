import React from "react";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";

interface FormInput {
  apiKey: string;
}

export default function OrganizationSettings() {
  const trpc = useTRPC();
  const upsertApiTokenMutation = useMutation(
    trpc.organization.upsertHetznerCredential.mutationOptions(),
  );
  const organization = useOrganizationContext();

  const {
    register,
    handleSubmit,
    reset,
    formState: { isDirty },
  } = useForm<FormInput>({
    defaultValues: { apiKey: "" },
  });

  const onSubmit = (data: FormInput) => {
    upsertApiTokenMutation.mutate(
      {
        id: organization.id,
        apiKey: data.apiKey,
      },
      {
        onSuccess: () => {
          reset();
        },
      },
    );
  };

  return (
    <div className="px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="text-xl font-bold">Organization Settings</h1>
      </div>
      <div className="w-full py-4 xl:w-3/5">
        <div className="border-mauve-6 rounded-md border-1 text-sm">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 border-b px-4 py-3 text-xs font-bold uppercase">
            General
          </div>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="flex items-center justify-between gap-2 px-4 py-2">
              <div>
                <h1 className="text-mauve-12">Hetzner API token</h1>
                <p className="text-mauve-11 text-xs">
                  Learn how to generate your API token{" "}
                  <a
                    className="text-mauve-11 text-xs underline"
                    target="_blank"
                    rel="noreferrer"
                    href="https://docs.hetzner.com/cloud/api/getting-started/generating-api-token/"
                  >
                    here
                  </a>
                  .
                </p>
              </div>
              <input
                className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-96 min-w-52 rounded-md border p-2 text-sm"
                type="text"
                placeholder="API Key*"
                {...register("apiKey")}
              />
            </div>
            {isDirty && (
              <div className="flex justify-end gap-1 px-4 pb-2">
                <Button
                  size="xs"
                  className="w-20"
                  intent="secondary"
                  onClick={() => reset()}
                >
                  Cancel
                </Button>
                <Button className="w-20" size="xs" type="submit">
                  Save
                </Button>
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
}

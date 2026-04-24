import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { Eye, EyeSlash } from "~/components/atoms/icons";

interface FormInput {
  apiKey: string;
}

export default function OrganizationGeneral() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();

  const upsertApiTokenMutation = useMutation(
    trpc.organization.upsertHetznerCredential.mutationOptions(),
  );
  const organization = useOrganizationContext();

  const { data: hetznerCredentialData, isLoading: isHetznerCredentialLoading } =
    useQuery(
      trpc.organization.getHetznerCredential.queryOptions({
        id: organization.id,
      }),
    );

  const {
    register,
    handleSubmit,
    reset,
    formState: { isDirty },
  } = useForm<FormInput>({
    defaultValues: { apiKey: "" },
  });

  useEffect(() => {
    reset({ apiKey: hetznerCredentialData?.credential?.secret ?? "" });
  }, [hetznerCredentialData?.credential?.secret, reset]);

  const onSubmit = (data: FormInput) => {
    upsertApiTokenMutation.mutate(
      {
        id: organization.id,
        apiKey: data.apiKey,
      },
      {
        onSuccess: async () => {
          reset({ apiKey: data.apiKey });
          await queryClient.invalidateQueries({
            queryKey: trpc.organization.getHetznerCredential.queryKey(),
          });
        },
      },
    );
  };

  const [show, setShow] = useState(false);

  if (!organization.isOwner) return null;

  return (
    <div className="w-full">
      <div className="border-mauve-6 rounded-md border-1 text-sm shadow-xs">
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
            {isHetznerCredentialLoading ? (
              <div className="flex w-1/2 items-center">
                <Skeleton className="mr-3 h-9.5 w-full" />
                <Eye className="h-4 w-4 flex-shrink-0" />
              </div>
            ) : (
              <div className="flex w-1/2 items-center">
                <input
                  className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border p-2 text-sm"
                  type={show ? "text" : "password"}
                  placeholder="API Key*"
                  {...register("apiKey")}
                />
                <button
                  onClick={() => setShow(!show)}
                  type="button"
                  className="pl-3"
                >
                  {show ? (
                    <EyeSlash className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </button>
              </div>
            )}
          </div>
          {isDirty && (
            <div className="flex justify-end gap-1 px-4 pb-2">
              <Button
                size="xs"
                className="w-20"
                intent="secondary"
                disabled={upsertApiTokenMutation.isPending}
                onClick={() => reset()}
              >
                Cancel
              </Button>
              <Button
                className="w-20"
                size="xs"
                type="submit"
                disabled={upsertApiTokenMutation.isPending}
              >
                Save
              </Button>
            </div>
          )}
        </form>
      </div>
    </div>
  );
}

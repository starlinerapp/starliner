import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Button from "~/components/atoms/button/Button";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { Eye, EyeSlash } from "~/components/atoms/icons";

interface FormInput {
  apiKey: string;
}

export default function HetznerCredentials() {
  const trpc = useTRPC();
  const queryClient = useQueryClient();
  const organization = useOrganizationContext();

  const upsertApiTokenMutation = useMutation(
    trpc.organization.upsertHetznerCredential.mutationOptions(),
  );

  const { data: hetznerCredentialData, isLoading } = useQuery(
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

  const [show, setShow] = useState(false);

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

  if (!organization.isOwner) return null;

  return (
    <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
      <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center border-b px-4 text-xs font-bold uppercase">
        Hetzner
      </div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="flex items-center justify-between gap-2 px-4 py-2">
          <div>
            <h2 className="text-mauve-12">API token</h2>
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
          {isLoading ? (
            <div className="relative w-1/2">
              <Skeleton className="h-9.5 w-full rounded-md" />
              <Eye className="text-mauve-11 absolute top-1/2 right-3 h-4 w-4 -translate-y-1/2" />
            </div>
          ) : (
            <div className="relative w-1/2">
              <input
                className="border-mauve-6 text-mauve-11 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border p-2 pr-10 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                type={show ? "text" : "password"}
                placeholder="API Key*"
                {...register("apiKey")}
              />
              <button
                onClick={() => setShow(!show)}
                type="button"
                aria-label={show ? "Hide API token" : "Show API token"}
                className="text-mauve-11 absolute top-1/2 right-3 -translate-y-1/2"
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
  );
}

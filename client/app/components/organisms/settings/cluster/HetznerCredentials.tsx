import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import Button from "~/components/atoms/button/Button";
import { Eye, EyeSlash } from "~/components/atoms/icons";
import Skeleton from "~/components/atoms/skeleton/Skeleton";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

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
    <div className="rounded-md border border-mauve-6 text-sm shadow-xs">
      <div className="flex h-14 items-center rounded-t-md border-mauve-6 border-b bg-gray-2 px-4 font-bold text-mauve-12 text-xs uppercase">
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
              <Eye className="absolute top-1/2 right-3 h-4 w-4 -translate-y-1/2 text-mauve-11" />
            </div>
          ) : (
            <div className="relative w-1/2">
              <input
                className="w-full min-w-52 rounded-md border border-mauve-6 bg-gray-2 p-2 pr-10 text-mauve-11 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)] placeholder:text-mauve-11"
                type={show ? "text" : "password"}
                placeholder="API Key*"
                {...register("apiKey")}
              />
              <button
                onClick={() => setShow(!show)}
                type="button"
                aria-label={show ? "Hide API token" : "Show API token"}
                className="absolute top-1/2 right-3 -translate-y-1/2 text-mauve-11"
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

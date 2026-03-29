import React from "react";
import { redirect, useLoaderData, useNavigate } from "react-router";
import type { Route } from "./+types/app";
import { auth } from "~/utils/auth/server";
import { ChevronDown, ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { type SubmitHandler, useForm } from "react-hook-form";

export async function loader({ request }: Route.LoaderArgs) {
  const url = new URL(request.url);

  const session = await auth.api.getSession({
    headers: request.headers,
  });

  if (!session) {
    const redirectTo = `${url.pathname}${url.search}`;
    return redirect(`/login?redirectTo=${encodeURIComponent(redirectTo)}`);
  }

  const installationIdParam = url.searchParams.get("installation_id");

  const installationId = installationIdParam
    ? Number(installationIdParam)
    : null;

  if (!installationId || Number.isNaN(installationId)) {
    throw new Response("Invalid installation_id", { status: 400 });
  }

  return {
    installationId,
  };
}

interface GithubAppFormInput {
  organizationId: number;
}

export default function GithubApp() {
  const { installationId } = useLoaderData<typeof loader>();

  const navigate = useNavigate();

  const trpc = useTRPC();
  const { data: organizations } = useQuery(
    trpc.organization.getUserOrganizations.queryOptions(),
  );
  const createGithubAppMutation = useMutation(
    trpc.githubApp.createGithubApp.mutationOptions(),
  );

  const { register, handleSubmit } = useForm<GithubAppFormInput>();

  const onSubmit: SubmitHandler<GithubAppFormInput> = async (data) => {
    createGithubAppMutation.mutate(
      {
        organizationId: Number(data.organizationId),
        installationId: installationId,
      },
      {
        onSuccess: () => {
          navigate("/");
        },
      },
    );
  };

  return (
    <div className="flex min-h-screen">
      <div className="bg-mauve-4 wiggle-pattern w-1/2"></div>
      <div className="flex w-1/2 items-center justify-center p-16 shadow-md">
        <div className="flex w-[500px] flex-col gap-4">
          <h1 className="text-xl font-medium">
            Link the Github App to an Organization
          </h1>
          <p className="text-mauve-11 text-sm">
            Linking the GitHub App enables GitHub features for this
            organization, such as listing repositories, reacting to pull
            requests, and automating workflows.
          </p>
          <form
            className="flex flex-col gap-2"
            onSubmit={handleSubmit(onSubmit)}
          >
            <span className="flex flex-col gap-1">
              <label htmlFor="name" className="text-sm">
                Organization
              </label>
              <div className="relative w-full">
                <select
                  {...register("organizationId", {
                    required: true,
                  })}
                  className="border-mauve-6 w-full appearance-none rounded-md border-1 p-2"
                >
                  {organizations?.map((org, i) => (
                    <option key={i} value={org.id}>
                      {org.name}
                    </option>
                  ))}
                </select>
                <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                  <ChevronDown width={15} className="stroke-mauve-10" />
                </div>
              </div>
            </span>
            <Button type="submit" className="mt-2" size="md">
              Complete Setup <ChevronRight className="w-4 stroke-3" />
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
}

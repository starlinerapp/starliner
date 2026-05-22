import { ArrowRight } from "~/components/atoms/icons";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import React from "react";
import { redirect, useLoaderData, useNavigate, useParams } from "react-router";
import { caller } from "~/utils/trpc/server";
import type { Route } from "./+types/githubapp";

export async function loader(loaderArgs: Route.LoaderArgs) {
  const { slug } = loaderArgs.params;

  if (!slug) {
    throw new Response("Not found", { status: 404 });
  }

  const trpcCaller = await caller(loaderArgs);
  const organizations = await trpcCaller.organization.getUserOrganizations();
  const organization = organizations?.find((o) => o.slug === slug);

  if (!organization) {
    throw new Response("Not found", { status: 404 });
  }

  const githubApp = await trpcCaller.githubApp.getGithubApp({
    organizationId: organization.id,
  });

  if (githubApp) {
    return redirect(`/${slug}`);
  }

  return {
    githubAppName: process.env.GITHUB_APP_NAME,
  };
}

export default function OrganizationGithubApp() {
  const { githubAppName } = useLoaderData<typeof loader>();
  const { slug } = useParams();
  const navigate = useNavigate();

  const redirectTo = `/${slug}`;

  return (
    <div className="flex w-125 flex-col gap-4">
      <button
        type="button"
        className="hover:bg-gray-4 flex cursor-pointer items-center gap-0.5 self-end rounded-md px-2 py-0.5 text-sm font-light"
        onClick={() => navigate(redirectTo)}
      >
        Skip <ArrowRight className="w-3" />
      </button>
      <h1 className="text-xl font-medium">Install the GitHub App</h1>
      <p className="text-mauve-11 text-sm">
        Connect GitHub to list repositories, react to pull requests, and
        automate workflows for your organization. You can install it later from
        Settings → Integrations.
      </p>
      <div className="mt-2">
        <InstallGitHubApp
          githubAppName={githubAppName}
          redirectTo={redirectTo}
        />
      </div>
    </div>
  );
}

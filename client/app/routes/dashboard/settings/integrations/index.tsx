import React from "react";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import { GithubLogo } from "~/components/atoms/icons";
import InstallGitHubApp from "~/components/atoms/github/InstallGitHubApp";
import { useLoaderData, useLocation } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { cn } from "~/utils/cn";
import ConfigureGitHubApp from "~/components/atoms/github/ConfigureGithubApp";
import Skeleton from "~/components/atoms/skeleton/Skeleton";

export function loader() {
  return {
    githubAppName: process.env.GITHUB_APP_NAME,
  };
}

export default function GitHubIntegration() {
  const { githubAppName } = useLoaderData<typeof loader>();
  const location = useLocation();

  const trpc = useTRPC();
  const organization = useOrganizationContext();

  const { data: githubApp, isLoading } = useQuery(
    trpc.githubApp.getGithubApp.queryOptions({
      organizationId: organization.id,
    }),
  );

  return (
    <>
      <Breadcrumbs
        crumbs={[
          { label: "Settings" },
          { label: "Integrations" },
          { label: "Integrations" },
        ]}
      />
      <div className="flex flex-col px-4 py-4">
        <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
          <div className="border-mauve-6 text-mauve-12 bg-gray-2 flex h-14 items-center border-b px-4 text-xs font-bold uppercase">
            Integrations
          </div>
          <form onSubmit={() => {}}>
            <div className="flex items-center justify-between gap-2 px-4 py-2">
              <div className="flex items-center gap-3">
                <GithubLogo className="h-7 w-7 invert" />
                <div>
                  <h2 className="text-mauve-12">GitHub App</h2>
                  {isLoading ? (
                    <div className="flex items-center gap-1.5">
                      <Skeleton className="h-2 w-2 rounded-full" />
                      <Skeleton className="h-5 w-20" />
                    </div>
                  ) : (
                    <div className="flex items-center gap-1.5">
                      <div
                        className={cn(
                          "h-2 w-2 rounded-full",
                          githubApp ? "bg-violet-10" : "bg-gray-10",
                        )}
                      ></div>
                      <p
                        className={cn(
                          "text-mauve-11 text-sm",
                          githubApp && "text-violet-11",
                        )}
                      >
                        {githubApp ? "Installed" : "Not installed"}
                      </p>
                    </div>
                  )}
                </div>
              </div>
              {isLoading ? (
                <Skeleton className="h-9 w-38 rounded-md" />
              ) : githubApp && location ? (
                <ConfigureGitHubApp
                  githubAppName={githubAppName}
                  redirectTo={location.pathname}
                />
              ) : (
                <InstallGitHubApp
                  githubAppName={githubAppName}
                  redirectTo={location.pathname}
                />
              )}
            </div>
          </form>
        </div>
      </div>
    </>
  );
}

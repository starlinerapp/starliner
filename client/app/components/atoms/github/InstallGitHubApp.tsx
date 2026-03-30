import React from "react";
import { GithubLogo } from "~/components/atoms/icons";

interface InstallGitHubAppProps {
  githubAppName: string | undefined;
}

export default function InstallGitHubApp({
  githubAppName,
}: InstallGitHubAppProps) {
  return (
    <a
      href={`https://github.com/apps/${githubAppName}/installations/new`}
      className="bg-gray-12 hover:bg-gray-11 flex w-32 cursor-pointer justify-center rounded-md px-3 py-2 text-sm text-white"
    >
      <span className="flex items-center gap-2.5">
        <GithubLogo className="h-5 w-5" />
        <p>Install App</p>
      </span>
    </a>
  );
}

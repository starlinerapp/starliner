import React from "react";
import { GithubLogo } from "~/components/atoms/icons";

export default function InstallGitHubApp() {
  return (
    <a
      // TODO: Move Github App Name to Env
      href="https://github.com/apps/dev-starliner-app/installations/new"
      className="bg-gray-12 hover:bg-gray-11 flex w-36 cursor-pointer justify-center rounded-md px-3 py-2 text-sm text-white"
    >
      <span className="flex items-center gap-3">
        <GithubLogo className="h-5 w-5" />
        <p>Install App</p>
      </span>
    </a>
  );
}

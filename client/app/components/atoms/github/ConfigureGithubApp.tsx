import { GithubLogo } from "~/components/atoms/icons";

interface ConfigureGitHubAppProps {
  githubAppName: string | undefined;
  redirectTo?: string;
}

export default function ConfigureGitHubApp({
  githubAppName,
  redirectTo,
}: ConfigureGitHubAppProps) {
  if (!githubAppName) {
    return null;
  }

  const installUrl = new URL(
    `https://github.com/apps/${githubAppName}/installations/new`,
  );
  if (redirectTo) {
    installUrl.searchParams.set("state", redirectTo);
  }

  return (
    <a
      href={installUrl.toString()}
      className="flex w-38 cursor-pointer justify-center rounded-md bg-gray-12 px-3 py-2 text-sm text-white hover:bg-gray-11"
    >
      <span className="flex items-center gap-2.5">
        <GithubLogo className="h-5 w-5" />
        <p>Configure App</p>
      </span>
    </a>
  );
}

import { GithubLogo } from "~/components/atoms/icons";

interface InstallGitHubAppProps {
  githubAppName: string | undefined;
  redirectTo?: string;
}

export default function InstallGitHubApp({
  githubAppName,
  redirectTo,
}: InstallGitHubAppProps) {
  const installUrl = new URL(
    `https://github.com/apps/${githubAppName}/installations/new`,
  );
  if (redirectTo) {
    installUrl.searchParams.set("state", redirectTo);
  }

  return (
    <a
      href={installUrl.toString()}
      className="flex w-32 cursor-pointer justify-center rounded-md bg-gray-12 px-3 py-2 text-sm text-white hover:bg-gray-11"
    >
      <span className="flex items-center gap-2.5">
        <GithubLogo className="h-5 w-5" />
        <p>Install App</p>
      </span>
    </a>
  );
}

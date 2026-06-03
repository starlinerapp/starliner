import { useSearchParams } from "react-router";
import { GithubLogo } from "~/components/atoms/icons";
import { getAuthClient } from "~/utils/auth/client";

export default function SignInWithGitHub() {
  const authClient = getAuthClient();

  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo") || "/";

  async function handleButtonClicked() {
    const callbackURL = new URL(redirectTo, window.location.origin).href;

    await authClient.signIn.social({
      provider: "github",
      callbackURL,
    });
  }

  return (
    <button
      type="button"
      onClick={handleButtonClicked}
      className="border-mauve-6 hover:bg-gray-3 text-mauve-12 flex w-full cursor-pointer justify-center rounded-md border bg-white px-4 py-2 text-sm"
    >
      <span className="flex items-center gap-2.5">
        <GithubLogo className="h-5 w-5 invert" />
        <p>Sign in with GitHub</p>
      </span>
    </button>
  );
}

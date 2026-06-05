import { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { NavLink, useSearchParams } from "react-router";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";
import SuccessBanner from "~/components/atoms/banner/SuccessBanner";
import Button from "~/components/atoms/button/Button";
import {
  ArrowRight,
  ChevronRight,
  Eye,
  EyeSlash,
} from "~/components/atoms/icons";
import { getAuthClient } from "~/utils/auth/client";

interface SignUpFormInput {
  email: string;
  password: string;
  username: string;
}

export default function SignUp() {
  const authClient = getAuthClient();
  const { register, handleSubmit, reset } = useForm<SignUpFormInput>();
  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo") || "/";

  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const onSubmit: SubmitHandler<SignUpFormInput> = async (data) => {
    const callbackURL = new URL(redirectTo, window.location.origin).href;

    await authClient.signUp.email(
      {
        email: data.email,
        password: data.password,
        name: data.username,
        callbackURL,
      },
      {
        onRequest: () => {
          // show loading state
        },
        onSuccess: () => {
          reset();
          setSuccess(
            "We sent you a verification email. Please verify your account before signing in.",
          );
        },
        onError: (ctx) => {
          setError(ctx.error.message);
        },
      },
    );
  };

  const loginLink =
    redirectTo !== "/"
      ? `/login?redirectTo=${encodeURIComponent(redirectTo)}`
      : "/login";

  return (
    <div className="flex w-125 flex-col gap-4">
      <p className="flex items-center justify-end gap-2 py-0.5 font-light text-sm">
        Already have an account?
        <NavLink
          to={loginLink}
          className="flex cursor-pointer items-center gap-1 rounded-md px-2 py-0.5 underline hover:bg-gray-4"
        >
          Sign in <ArrowRight className="w-3" />
        </NavLink>
      </p>
      <h1 className="font-medium text-xl">Sign up for Starliner</h1>
      {error && <ErrorBanner text={error} />}
      {success && <SuccessBanner text={success} />}
      <form className="flex flex-col gap-2" onSubmit={handleSubmit(onSubmit)}>
        <span className="flex flex-col gap-1">
          <label htmlFor="username" className="text-sm">
            Full Name
          </label>
          <input
            id="username"
            className="rounded-md border border-mauve-6 p-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
            type="text"
            placeholder="Full Name*"
            {...register("username")}
          />
        </span>
        <span className="flex flex-col gap-1">
          <label htmlFor="email" className="text-sm">
            Email
          </label>
          <input
            id="email"
            className="rounded-md border border-mauve-6 p-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
            type="email"
            placeholder="Email"
            {...register("email")}
          />
        </span>
        <span className="flex flex-col gap-1">
          <label htmlFor="password" className="text-sm">
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              className="w-full rounded-md border border-mauve-6 p-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              {...register("password")}
            />

            <button
              onClick={() => setShowPassword(!showPassword)}
              type="button"
              aria-label={showPassword ? "Hide password" : "Show password"}
              className="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer text-mauve-11"
            >
              {showPassword ? (
                <EyeSlash className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </button>
          </div>
        </span>
        <Button className="mt-2" type="submit" size="md">
          Create account <ChevronRight className="w-4 stroke-3" />
        </Button>
        <p className="mt-1 text-mauve-11 text-xs">
          By creating an account you agree to the Terms of Service. For more
          information about Starliner&#39;s privacy practices see the Privacy
          Statement. We&#39;ll occasionally send account-related emails.
        </p>
      </form>
    </div>
  );
}

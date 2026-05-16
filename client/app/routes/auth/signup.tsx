import React, { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import { NavLink, useNavigate, useSearchParams } from "react-router";
import Button from "~/components/atoms/button/Button";
import { getAuthClient } from "~/utils/auth/client";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

interface SignUpFormInput {
  email: string;
  password: string;
  username: string;
}

export default function SignUp() {
  const authClient = getAuthClient();
  const { register, handleSubmit } = useForm<SignUpFormInput>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo") || "/";

  const [error, setError] = useState<string | null>(null);

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
          navigate(redirectTo);
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
    <div className="flex w-[500px] flex-col gap-4">
      <p className="flex items-center justify-end gap-2 py-0.5 text-sm font-light">
        Already have an account?
        <NavLink
          to={loginLink}
          className="hover:bg-gray-4 flex cursor-pointer items-center gap-1 rounded-md px-2 py-0.5 underline"
        >
          Sign in <ArrowRight className="w-3" />
        </NavLink>
      </p>
      <h1 className="text-xl font-medium">Sign up for Starliner</h1>
      {error && <ErrorBanner text={error} />}
      <form className="flex flex-col gap-2" onSubmit={handleSubmit(onSubmit)}>
        <span className="flex flex-col gap-1">
          <label htmlFor="email" className="text-sm">
            Email
          </label>
          <input
            className="border-mauve-6 rounded-md border-1 p-2"
            type="text"
            placeholder="Email"
            {...register("email")}
          />
        </span>
        <span className="flex flex-col gap-1">
          <label htmlFor="password" className="text-sm">
            Password
          </label>
          <input
            className="border-mauve-6 rounded-md border-1 p-2"
            type="password"
            placeholder="Password"
            {...register("password")}
          />
        </span>
        <span className="flex flex-col gap-1">
          <label htmlFor="username" className="text-sm">
            Username
          </label>
          <input
            className="border-mauve-6 rounded-md border-1 p-2"
            type="text"
            placeholder="Username"
            {...register("username")}
          />
        </span>
        <Button className="mt-2" type="submit" size="md">
          Create account <ChevronRight className="w-4 stroke-3" />
        </Button>
        <p className="text-mauve-11 mt-1 text-xs">
          By creating an account you agree to the Terms of Service. For more
          information about Starliner&#39;s privacy practices see the Privacy
          Statement. We&#39;ll occasionally send account-related emails.
        </p>
      </form>
    </div>
  );
}

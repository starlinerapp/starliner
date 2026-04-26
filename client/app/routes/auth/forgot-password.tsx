import React, { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import { NavLink, useSearchParams } from "react-router";
import Button from "~/components/atoms/button/Button";
import { authClient } from "~/utils/auth/client";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

interface ForgotFormInput {
  email: string;
}

export default function ForgotPassword() {
  const { register, handleSubmit } = useForm<ForgotFormInput>();
  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo") || "/";

  const [error, setError] = useState<string | null>(null);
  const [submitted, setSubmitted] = useState(false);

  const onSubmit: SubmitHandler<ForgotFormInput> = async (data) => {
    const redirectToReset = new URL("/reset-password", window.location.origin);
    if (redirectTo && redirectTo !== "/") {
      redirectToReset.searchParams.set("redirectTo", redirectTo);
    }

    await authClient.requestPasswordReset(
      {
        email: data.email,
        redirectTo: redirectToReset.toString(),
      },
      {
        onRequest: () => setError(null),
        onSuccess: () => setSubmitted(true),
        onError: (ctx) => setError(ctx.error.message),
      },
    );
  };

  const loginLink =
    redirectTo !== "/"
      ? `/login?redirectTo=${encodeURIComponent(redirectTo)}`
      : "/login";

  return (
    <div className="flex w-[500px] flex-col gap-4">
      <p className="flex items-center justify-end gap-1.5 py-0.5 text-sm font-light">
        Remember your password?
        <NavLink
          to={loginLink}
          className="hover:bg-gray-4 flex cursor-pointer items-center gap-1 rounded-md px-2 py-0.5 underline"
        >
          Sign in <ArrowRight className="w-3" />
        </NavLink>
      </p>
      <h1 className="text-xl font-medium">Reset your password</h1>
      <p className="text-mauve-11 -mt-1 text-sm">
        Enter your account email and we’ll send you a link to set a new
        password.
      </p>
      {error && <ErrorBanner text={error} />}
      {submitted ? (
        <p className="text-mauve-12">
          If an account exists for that address, you’ll get an email with reset
          instructions. You can close this page.
        </p>
      ) : (
        <form className="flex flex-col gap-2" onSubmit={handleSubmit(onSubmit)}>
          <span className="flex flex-col gap-1">
            <label htmlFor="email" className="text-sm">
              Email
            </label>
            <input
              className="border-mauve-6 rounded-md border-1 p-2"
              type="email"
              id="email"
              autoComplete="email"
              placeholder="Email"
              {...register("email", { required: true })}
            />
          </span>
          <Button className="mt-2" type="submit" size="md">
            Send reset link <ChevronRight className="w-4 stroke-3" />
          </Button>
        </form>
      )}
    </div>
  );
}

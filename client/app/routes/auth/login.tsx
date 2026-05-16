import React, { useCallback, useState } from "react";
import { NavLink, useNavigate, useSearchParams } from "react-router";
import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { type SubmitHandler, useForm } from "react-hook-form";
import { getAuthClient } from "~/utils/auth/client";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

const VERIFY_EMAIL_ERROR = "Please verify your email address";

interface LoginFormInput {
  email: string;
  password: string;
}

export default function Login() {
  const authClient = getAuthClient();
  const { register, handleSubmit } = useForm<LoginFormInput>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo") || "/";

  const [error, setError] = useState<string | null>(null);
  const [unverifiedEmail, setUnverifiedEmail] = useState<string | null>(null);
  const [isResending, setIsResending] = useState(false);

  const handleResendVerification = useCallback(async () => {
    if (!unverifiedEmail) {
      return;
    }
    setIsResending(true);
    try {
      const { error: resendError } = await authClient.sendVerificationEmail({
        email: unverifiedEmail,
        callbackURL: new URL(redirectTo, window.location.origin).href,
      });
      if (resendError) {
        setError(
          "message" in resendError && resendError.message
            ? String(resendError.message)
            : "Could not send the verification email.",
        );
        return;
      }
      setError(null);
      setUnverifiedEmail(null);
    } finally {
      setIsResending(false);
    }
  }, [authClient, unverifiedEmail, redirectTo]);

  const onSubmit: SubmitHandler<LoginFormInput> = async (data) => {
    await authClient.signIn.email(
      {
        email: data.email,
        password: data.password,
      },
      {
        onSuccess: () => {
          navigate(redirectTo);
        },
        onError: (ctx) => {
          if (ctx.error.status === 403) {
            setError(VERIFY_EMAIL_ERROR);
            setUnverifiedEmail(data.email);
            return;
          }
          setUnverifiedEmail(null);
          setError(ctx.error.message);
        },
      },
    );
  };

  const signupLink =
    redirectTo !== "/"
      ? `/signup?redirectTo=${encodeURIComponent(redirectTo)}`
      : "/signup";

  return (
    <div className="flex w-[500px] flex-col gap-4">
      <p className="flex items-center justify-end gap-1.5 py-0.5 text-sm font-light">
        Don&#39;t have an account?
        <NavLink
          to={signupLink}
          className="hover:bg-gray-4 flex cursor-pointer items-center gap-1 rounded-md px-2 py-0.5 underline"
        >
          Sign up <ArrowRight className="w-3" />
        </NavLink>
      </p>
      <h1 className="text-xl font-medium">Sign in to Starliner</h1>
      {error && (
        <ErrorBanner text={error}>
          {unverifiedEmail && (
            <button
              className="enabled:hover:text-mauve-12 text-sm underline enabled:cursor-pointer"
              type="button"
              disabled={isResending}
              onClick={() => {
                void handleResendVerification();
              }}
            >
              {isResending ? "Sending…" : "Resend verification email"}
            </button>
          )}
        </ErrorBanner>
      )}
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
          <div className="flex justify-between text-sm">
            <label htmlFor="password">Password</label>
            <NavLink
              to={
                redirectTo !== "/"
                  ? `/forgot-password?redirectTo=${encodeURIComponent(redirectTo)}`
                  : "/forgot-password"
              }
              className="text-mauve-11 hover:text-mauve-12"
            >
              Forgot password?
            </NavLink>
          </div>

          <input
            className="border-mauve-6 rounded-md border-1 p-2"
            type="password"
            placeholder="Password"
            {...register("password")}
          />
        </span>
        <Button className="mt-2" size="md" type="submit">
          Sign in <ChevronRight className="w-4 stroke-3" />
        </Button>
      </form>
    </div>
  );
}

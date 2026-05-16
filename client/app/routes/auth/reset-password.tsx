import React, { useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { NavLink, useNavigate, useSearchParams } from "react-router";
import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { getAuthClient } from "~/utils/auth/client";
import ErrorBanner from "~/components/atoms/banner/ErrorBanner";

interface ResetFormInput {
  password: string;
  confirmPassword: string;
}

const INVALID_TOKEN_ERROR =
  "This reset link is invalid or has expired. Request a new one from the sign-in page.";

export default function ResetPassword() {
  const authClient = getAuthClient();
  const { register, handleSubmit } = useForm<ResetFormInput>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const token = searchParams.get("token");
  const queryError = searchParams.get("error");
  const redirectTo = searchParams.get("redirectTo") || "/";

  const [error, setError] = useState<string | null>(null);

  const isInvalidRequest =
    queryError === "INVALID_TOKEN" || token == null || token === "";

  const onSubmit: SubmitHandler<ResetFormInput> = async (data) => {
    if (!token) {
      setError(INVALID_TOKEN_ERROR);
      return;
    }
    if (data.password !== data.confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    await authClient.resetPassword(
      {
        newPassword: data.password,
        token,
      },
      {
        onRequest: () => setError(null),
        onSuccess: () => {
          const href =
            redirectTo !== "/"
              ? `/login?redirectTo=${encodeURIComponent(redirectTo)}`
              : "/login";
          navigate(href, { replace: true });
        },
        onError: (ctx) => setError(ctx.error.message),
      },
    );
  };

  const loginLink =
    redirectTo !== "/"
      ? `/login?redirectTo=${encodeURIComponent(redirectTo)}`
      : "/login";
  const forgotLink =
    redirectTo !== "/"
      ? `/forgot-password?redirectTo=${encodeURIComponent(redirectTo)}`
      : "/forgot-password";

  if (isInvalidRequest) {
    return (
      <div className="flex w-[500px] flex-col gap-4">
        <h1 className="text-xl font-medium">Oups! The link is invalid</h1>
        <p className="text-mauve-11 text-sm">{INVALID_TOKEN_ERROR}</p>
        <div className="hover:bg-gray-3 flex w-56 items-center gap-2 rounded-md p-1 py-0.5">
          <ArrowRight width="20" strokeWidth="2" />
          <NavLink to={forgotLink} className="underline">
            Request a new reset link
          </NavLink>
        </div>
        <div className="hover:bg-gray-3 flex w-38 items-center gap-2 rounded-md p-1 py-0.5">
          <ArrowRight width="20" strokeWidth="2" />
          <NavLink to={loginLink} className="underline">
            Back to sign in
          </NavLink>
        </div>
      </div>
    );
  }

  return (
    <div className="flex w-[500px] flex-col gap-4">
      <p className="flex items-center justify-end text-sm font-light">
        <NavLink
          to={loginLink}
          className="hover:bg-gray-4 flex cursor-pointer items-center gap-1 rounded-md px-2 py-0.5 underline"
        >
          Back to sign in <ArrowRight className="w-3" />
        </NavLink>
      </p>
      <h1 className="text-xl font-medium">Set a new password</h1>
      <p className="text-mauve-11 -mt-1 text-sm">
        Choose a new password for your account.
      </p>
      {error && <ErrorBanner text={error} />}
      <form className="flex flex-col gap-2" onSubmit={handleSubmit(onSubmit)}>
        <span className="flex flex-col gap-1">
          <label htmlFor="password" className="text-sm">
            New password
          </label>
          <input
            className="border-mauve-6 rounded-md border-1 p-2"
            type="password"
            id="password"
            autoComplete="new-password"
            placeholder="New password"
            {...register("password", { required: true, minLength: 8 })}
          />
        </span>
        <span className="flex flex-col gap-1">
          <label htmlFor="confirmPassword" className="text-sm">
            Confirm password
          </label>
          <input
            className="border-mauve-6 rounded-md border-1 p-2"
            type="password"
            id="confirmPassword"
            autoComplete="new-password"
            placeholder="Confirm new password"
            {...register("confirmPassword", { required: true })}
          />
        </span>
        <Button className="mt-2" type="submit" size="md">
          Update password <ChevronRight className="w-4 stroke-3" />
        </Button>
      </form>
    </div>
  );
}

import React from "react";
import { NavLink, useNavigate } from "react-router";
import { ArrowRight, ChevronRight } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { type SubmitHandler, useForm } from "react-hook-form";
import { authClient } from "~/utils/auth/client";

interface LoginFormInput {
  email: string;
  password: string;
}

export default function Login() {
  const { register, handleSubmit } = useForm<LoginFormInput>();
  const navigate = useNavigate();

  const onSubmit: SubmitHandler<LoginFormInput> = async (data) => {
    await authClient.signIn.email(
      {
        email: data.email,
        password: data.password,
      },
      {
        onRequest: () => {
          // show loading state
        },
        onSuccess: () => {
          navigate("/");
        },
        onError: (ctx) => {
          console.log(ctx.error);
        },
      },
    );
  };

  return (
    <div className="flex w-[500px] flex-col gap-4">
      <p className="flex items-center justify-end gap-1.5 py-0.5 text-sm font-light">
        Don&#39;t have an account?
        <NavLink
          to="/signup"
          className="hover:bg-gray-4 flex cursor-pointer items-center gap-1 rounded-md px-2 py-0.5 underline"
        >
          Sign up <ArrowRight className="w-3" />
        </NavLink>
      </p>
      <h1 className="text-xl font-medium">Sign in to Starliner</h1>
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
        <Button type="submit">
          Sign in <ChevronRight className="w-4 stroke-3" />
        </Button>
      </form>
    </div>
  );
}

import React, {
  type ButtonHTMLAttributes,
  type PropsWithChildren,
} from "react";

type ButtonProps = PropsWithChildren<ButtonHTMLAttributes<HTMLButtonElement>>;

export default function Button({ children, ...props }: ButtonProps) {
  return (
    <button
      className="bg-violet-10 hover:bg-violet-9 text-white-a12 mt-4 flex w-full cursor-pointer items-center justify-center gap-2 rounded-md px-4 py-3 text-center"
      {...props}
    >
      {children}
    </button>
  );
}

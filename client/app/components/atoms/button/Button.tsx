import React, {
  type ButtonHTMLAttributes,
  type PropsWithChildren,
} from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "~/utils/cn";

type ButtonProps = PropsWithChildren<ButtonHTMLAttributes<HTMLButtonElement>>;

const buttonVariants = cva(
  [
    "flex",
    "w-full",
    "cursor-pointer",
    "items-center",
    "justify-center",
    "gap-2",
    "rounded-md",
    "text-center",
    "text-sm",
    "border-1",
  ],
  {
    variants: {
      size: { xs: "p-1.5", sm: "p-2", md: "px-4 py-3" },
      intent: {
        primary: "bg-violet-10 hover:bg-violet-9 text-white",
        secondary: "bg-white border-mauve-7 hover:bg-mauve-1 text-mauve-12",
        danger:
          "text-red-11 font-bold border-gray-7 bg-gray-3 hover:text-white hover:bg-red-11 hover:border-red-8",
        text: "p-1 text-xs border-0 justify-start font-medium p-0 gap-1 text-violet-11",
      },
      disabled: {
        false: null,
        true: [
          "text-mauve-11",
          "hover:text-mauve-11",
          "bg-white",
          "hover:bg-white",
          "border-mauve-6",
          "hover:border-mauve-6",
          "cursor-not-allowed",
        ],
      },
    },
    defaultVariants: {
      size: "sm",
      intent: "primary",
      disabled: false,
    },
  },
);

export interface ButtonVariants
  extends Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, "disabled">,
    VariantProps<typeof buttonVariants> {}

export default function Button({
  children,
  className,
  size,
  intent,
  disabled,
  ...props
}: ButtonProps & ButtonVariants) {
  return (
    <button
      className={cn(buttonVariants({ intent, size, disabled }), className)}
      {...props}
    >
      {children}
    </button>
  );
}

import { cva, type VariantProps } from "class-variance-authority";
import { CheckCircle } from "lucide-react";
import type React from "react";
import { Link } from "react-router";
import { ExclamationTriangle, LinkOut } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

const bannerVariants = cva(["flex", "w-full", "rounded-md", "border-1"], {
  variants: {
    intent: {
      error: "border-red-6 bg-red-3",
      success: "border-grass-6 bg-grass-3",
      warning: "border-amber-6 bg-amber-3",
    },
  },
  defaultVariants: {
    intent: "error",
  },
});

const bannerIconVariants = cva(
  ["flex", "w-11", "items-center", "justify-center", "rounded-l-sm"],
  {
    variants: {
      intent: {
        error: "bg-red-9",
        success: "bg-grass-9",
        warning: "bg-amber-9",
      },
    },
    defaultVariants: {
      intent: "error",
    },
  },
);

const bannerIcons = {
  error: <ExclamationTriangle width={18} strokeWidth={2} />,
  success: <CheckCircle width={18} strokeWidth={2} />,
  warning: <ExclamationTriangle width={18} strokeWidth={2} />,
};

interface BannerLinkOutProps {
  text: string;
  href: string;
}

export interface BannerVariants extends VariantProps<typeof bannerVariants> {}

interface BannerProps extends BannerVariants {
  text: string;
  description?: string;
  linkOut?: BannerLinkOutProps;
  children?: React.ReactNode;
  className?: string;
}

export default function Banner({
  text,
  description,
  linkOut,
  children,
  className,
  intent,
}: BannerProps) {
  return (
    <div className={cn(bannerVariants({ intent }), className)}>
      <div className={bannerIconVariants({ intent })}>
        {bannerIcons[intent ?? "error"]}
      </div>
      <div className="flex grow items-center gap-2 p-2.5">
        <div className="flex flex-col">
          <p className="font-light text-sm">{text}</p>
          {description && (
            <p className="font-light text-mauve-11 text-xs">{description}</p>
          )}
        </div>
        {children}
        {linkOut && (
          <span className="flex items-center gap-1">
            <Link className="font-light text-sm underline" to={linkOut.href}>
              {linkOut.text}
            </Link>
            <LinkOut width={18} />
          </span>
        )}
      </div>
    </div>
  );
}

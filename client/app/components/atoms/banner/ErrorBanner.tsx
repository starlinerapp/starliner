import type React from "react";
import { Link } from "react-router";
import { ExclamationTriangle, LinkOut } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

interface ErrorBannerLinkOutProps {
  text: string;
  href: string;
}

interface ErrorBannerProps {
  text: string;
  linkOut?: ErrorBannerLinkOutProps;
  children?: React.ReactNode;
  className?: string;
}

export default function ErrorBanner({
  text,
  linkOut,
  children,
  className,
}: ErrorBannerProps) {
  return (
    <div
      className={cn(
        "flex w-full rounded-md border-1 border-red-6 bg-red-3",
        className,
      )}
    >
      <div className="flex w-11 items-center justify-center rounded-l-sm bg-red-9">
        <ExclamationTriangle width={18} strokeWidth={2} />
      </div>
      <div className="flex items-center gap-2 p-2.5">
        <p className="font-light text-sm">{text}</p>
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

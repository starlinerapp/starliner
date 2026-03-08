import React from "react";
import { cn } from "~/utils/cn";
import { ExclamationTriangle, LinkOut } from "~/components/atoms/icons";
import { Link } from "react-router";

interface ErrorBAnnerLinkOutProps {
  text: string;
  href: string;
}

interface ErrorBannerProps {
  text: string;
  linkOut?: ErrorBAnnerLinkOutProps;
  className?: string;
}

export default function ErrorBanner({
  text,
  linkOut,
  className,
}: ErrorBannerProps) {
  return (
    <div
      className={cn(
        "border-red-6 bg-red-3 flex w-full rounded-md border-1",
        className,
      )}
    >
      <div className="bg-red-9 flex w-11 items-center justify-center rounded-l-md">
        <ExclamationTriangle width={18} strokeWidth={2} />
      </div>
      <div className="flex items-center gap-2 p-2.5">
        <p className="text-sm font-light">{text}</p>
        {linkOut && (
          <span className="flex items-center gap-1">
            <Link className="text-sm font-light underline" to={linkOut.href}>
              {linkOut.text}
            </Link>
            <LinkOut width={18} />
          </span>
        )}
      </div>
    </div>
  );
}

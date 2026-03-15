import React from "react";
import { ExclamationTriangle, LinkOut } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";
import { Link } from "react-router";

interface WarningBannerLinkOutProps {
  text: string;
  href: string;
}

interface WarningBannerProps {
  text: string;
  linkOut?: WarningBannerLinkOutProps;
  className?: string;
}

export default function WarningBanner({
  text,
  linkOut,
  className,
}: WarningBannerProps) {
  return (
    <div
      className={cn(
        "border-amber-6 bg-amber-3 flex w-full rounded-md border-1",
        className,
      )}
    >
      <div className="bg-amber-9 flex w-11 items-center justify-center rounded-l-sm">
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

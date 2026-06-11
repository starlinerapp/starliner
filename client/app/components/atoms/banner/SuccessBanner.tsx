import type React from "react";
import { Link } from "react-router";
import {CheckCircle, LinkOut} from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

interface SuccessBannerLinkOutProps {
  text: string;
  href: string;
}

interface SuccessBannerProps {
  text: string;
  linkOut?: SuccessBannerLinkOutProps;
  children?: React.ReactNode;
  className?: string;
}

export default function SuccessBanner({
  text,
  linkOut,
  children,
  className,
}: SuccessBannerProps) {
  return (
    <div
      className={cn(
        "flex w-full rounded-md border border-grass-6 bg-grass-3",
        className,
      )}
    >
      <div className="flex w-11 items-center justify-center rounded-l-sm bg-grass-9">
        <CheckCircle width={18} strokeWidth={2} />
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

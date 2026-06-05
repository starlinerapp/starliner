import { Link } from "react-router";
import { ExclamationTriangle, LinkOut } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

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
        "flex w-full rounded-md border-1 border-amber-6 bg-amber-3",
        className,
      )}
    >
      <div className="flex w-11 items-center justify-center rounded-l-sm bg-amber-9">
        <ExclamationTriangle width={18} strokeWidth={2} />
      </div>
      <div className="flex items-center gap-2 p-2.5">
        <p className="font-light text-sm">{text}</p>
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

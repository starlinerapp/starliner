import React, { useMemo } from "react";
import * as Popover from "@radix-ui/react-popover";
import * as HoverCard from "@radix-ui/react-hover-card";
import { cn } from "~/utils/cn";
import { ChevronRight, Plus } from "~/components/atoms/icons";
import { Link, useParams } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";

interface OrganizationIconProps {
  name: string;
  className?: string;
}

function OrganizationIcon({ name, className }: OrganizationIconProps) {
  return (
    <div
      className={cn(
        "bg-violet-9 flex h-10 w-10 cursor-pointer items-center justify-center rounded-md border p-1 text-lg text-white",
        className,
      )}
    >
      {name.substring(0, 1).toUpperCase()}
    </div>
  );
}

export default function OrganizationBadge() {
  const { slug } = useParams();
  const [open, setOpen] = React.useState(false);

  const trpc = useTRPC();
  const { data: organizations } = useQuery(
    trpc.organization.getUserOrganizations.queryOptions(),
  );

  const currentOrganization = organizations?.find((o) => o.slug === slug);
  const otherOrganizations = useMemo(() => {
    return (
      organizations?.filter((organization) => organization.slug !== slug) ?? []
    );
  }, [organizations, slug]);

  const { data: projects } = useQuery(
    trpc.organization.getUserProjects.queryOptions(
      { id: currentOrganization?.id ?? 0 },
      { enabled: !!currentOrganization?.id },
    ),
  );

  return (
    <Popover.Root open={open} onOpenChange={setOpen}>
      <Popover.Trigger className="data-[state=open]:bg-violet-3 data-[state=open]:border-gray-4 hover:bg-violet-3 hover:border-gray-4 flex h-11 w-11 items-center justify-center self-center rounded-md border border-white">
        <OrganizationIcon name={currentOrganization?.name ?? ""} />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          side="right"
          align="start"
          className="border-gray-6 m-2 rounded-md border bg-white shadow-md"
        >
          <div className="flex min-w-[180px] flex-col p-1">
            <div className="flex gap-2.5 p-1">
              <OrganizationIcon
                name={currentOrganization?.name ?? ""}
                className="h-9 w-9"
              />
              <div className="flex flex-col justify-between">
                <p className="text-gray-12 text-sm font-bold">
                  {currentOrganization?.name ?? ""}
                </p>
                <p className="text-gray-11 text-xs">
                  {projects?.length ?? 0}{" "}
                  {projects?.length === 1 ? "Project" : "Projects"}
                </p>
              </div>
            </div>
            <Link
              to={`/${currentOrganization?.slug}/settings/organization`}
              className="hover:bg-gray-3 flex flex-row items-center gap-2 rounded-md p-2 text-xs"
              onClick={() => setOpen(false)}
            >
              <p>Organization Settings</p>
            </Link>
            <Link
              to={`/${currentOrganization?.slug}/projects`}
              className="hover:bg-gray-3 flex flex-row items-center gap-2 rounded-md p-2 text-xs"
              onClick={() => setOpen(false)}
            >
              <p>Projects</p>
            </Link>
            <HoverCard.Root openDelay={0} closeDelay={100}>
              <HoverCard.Trigger className="hover:bg-gray-3 flex cursor-pointer flex-row items-center justify-between gap-2 rounded-md p-2 text-xs">
                <p>Switch Organization</p>
                <ChevronRight width={12} strokeWidth={2.5} />
              </HoverCard.Trigger>
              <HoverCard.Portal>
                <HoverCard.Content
                  side="right"
                  align="start"
                  sideOffset={-8}
                  alignOffset={-8}
                  className="border-gray-6 m-2 cursor-pointer rounded-md border bg-white shadow-md"
                >
                  <div className="flex min-w-[160px] flex-col p-1">
                    {otherOrganizations.map((organization, i) => (
                      <Link
                        to={`/${organization.slug}`}
                        key={i}
                        className="hover:bg-gray-3 flex items-center gap-2 rounded-md p-2 text-xs"
                        onClick={() => setOpen(false)}
                      >
                        {organization.name}
                      </Link>
                    ))}
                    {otherOrganizations.length > 0 && (
                      <hr className="border-gray-4 my-1" />
                    )}
                    <Link
                      className="hover:bg-gray-3 flex flex-row items-center justify-between gap-2 rounded-md p-2 text-xs"
                      to="/organizations/new"
                      target="_blank"
                      onClick={() => setOpen(false)}
                    >
                      <p className="flex items-center gap-1 text-xs">
                        <Plus width={15} strokeWidth={2} /> New organization
                      </p>
                    </Link>
                  </div>
                </HoverCard.Content>
              </HoverCard.Portal>
            </HoverCard.Root>
          </div>
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  );
}

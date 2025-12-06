import React from "react";
import SidebarItem from "~/components/molecules/sidebar/SidebarItem";
import Avatar from "~/components/atoms/avatar/Avatar";
import OrganizationBadge from "~/components/molecules/sidebar/OrganizationBadge";

type SidebarItem = {
  id: string;
  title: string;
  icon: React.ReactNode;
  href: string;
};

interface SidebarProps {
  sidebarItems: SidebarItem[];
  children: React.ReactNode;
}

export default function Sidebar({ sidebarItems, children }: SidebarProps) {
  return (
    <div className="flex">
      <div className="border-mauve-6 flex h-screen w-18 flex-col justify-between border-r-1 pt-4">
        <div className="flex flex-col gap-3 self-center">
          <OrganizationBadge />
          {sidebarItems.map((item) => (
            <SidebarItem
              key={item.href}
              title={item.title}
              icon={item.icon}
              href={item.href}
            />
          ))}
        </div>
        <div className="flex w-2/3 justify-center self-center pt-3 pb-2">
          <div className="flex h-11 w-11 items-center justify-center">
            <Avatar />
          </div>
        </div>
      </div>
      {children}
    </div>
  );
}

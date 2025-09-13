import React from "react";
import { Cog, InboxStack } from "~/components/atoms/icons";
import SidebarItem from "~/components/molecules/sidebar/SidebarItem";
import Avatar from "~/components/atoms/avatar/Avatar";

export default function Sidebar() {
  const sidebarItems = [
    {
      title: "Projects",
      icon: <InboxStack />,
      href: "/",
    },
    {
      title: "Settings",
      icon: <Cog />,
      href: "/settings",
    },
  ];

  return (
    <div className="border-mauve-6 flex h-screen w-18 flex-col justify-between border-r-1 pt-4">
      <div className="flex flex-col gap-2 self-center">
        <div className="bg-violet-9 mb-2 flex h-10 w-10 cursor-pointer items-center justify-center self-center rounded-md border-1 stroke-1 p-1 text-lg text-white">
          S
        </div>
        {sidebarItems.map((item) => (
          <SidebarItem
            key={item.href}
            title={item.title}
            icon={item.icon}
            href={item.href}
          />
        ))}
      </div>
      <div className="border-mauve-6 flex w-2/3 justify-center self-center border-t-1 pt-3 pb-2">
        <div className="flex h-10 w-10 items-center justify-center">
          <Avatar />
        </div>
      </div>
    </div>
  );
}

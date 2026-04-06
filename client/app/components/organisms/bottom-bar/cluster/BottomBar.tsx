import React, { useState } from "react";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";
import TerminalClient from "~/components/atoms/terminal/Terminal.client";

const navigationItems = ["Terminal"] as const;
type NavigationItem = (typeof navigationItems)[number];

interface BottomBarProps {
  clusterId: number;
}

export default function BottomBar({ clusterId }: BottomBarProps) {
  const [selected, setSelected] = useState<NavigationItem>("Terminal");

  return (
    <div className="-mt-1 flex h-full flex-col bg-white">
      <NavigationBar
        items={navigationItems}
        selected={selected}
        onSelect={setSelected}
      />
      {selected === "Terminal" && (
        <TerminalClient
          webSocketUrl={`wss://${window.location.host}/ws/clusters/${clusterId}`}
        />
      )}
    </div>
  );
}

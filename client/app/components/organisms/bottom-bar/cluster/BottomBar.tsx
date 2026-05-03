import React, { useState } from "react";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";
import TerminalClient from "~/components/atoms/terminal/Terminal.client";
import Logs from "~/components/organisms/bottom-bar/cluster/Logs";

const navigationItems = ["Logs", "Terminal"] as const;
type NavigationItem = (typeof navigationItems)[number];

interface BottomBarProps {
  clusterId: number;
  status: string;
}

export default function BottomBar({ clusterId, status }: BottomBarProps) {
  const [selected, setSelected] = useState<NavigationItem>(
    status === "running" ? "Terminal" : "Logs",
  );

  return (
    <div className="-mt-1 flex h-full flex-col bg-white">
      <NavigationBar
        items={navigationItems}
        selected={selected}
        onSelect={setSelected}
      />
      {selected === "Logs" ? (
        <div className="flex min-h-0 flex-1 flex-col p-4">
          <Logs clusterId={clusterId} />
        </div>
      ) : status === "running" ? (
        <TerminalClient
          webSocketUrl={`wss://${window.location.host}/ws/clusters/${clusterId}`}
        />
      ) : (
        <p className="text-mauve-11 p-4">
          Terminal is available once the cluster is running.
        </p>
      )}
    </div>
  );
}

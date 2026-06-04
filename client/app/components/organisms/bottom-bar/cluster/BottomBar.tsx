import { useState } from "react";
import TerminalClient from "~/components/atoms/terminal/Terminal.client";
import Logs from "~/components/organisms/bottom-bar/cluster/Logs";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";

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
        <div className="flex min-h-0 flex-1 flex-col">
          <TerminalClient
            webSocketUrl={`wss://${window.location.host}/ws/clusters/${clusterId}`}
          />
        </div>
      ) : (
        <p className="p-4 text-mauve-11">
          Terminal is available once the cluster is running.
        </p>
      )}
    </div>
  );
}

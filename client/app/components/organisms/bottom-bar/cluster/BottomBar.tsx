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
        <div
          key="logs"
          className="fade-in-0 zoom-in-95 flex min-h-0 flex-1 animate-in flex-col p-4 duration-200"
        >
          <Logs clusterId={clusterId} />
        </div>
      ) : status === "running" ? (
        <div
          key="terminal"
          className="fade-in-0 zoom-in-95 flex min-h-0 flex-1 animate-in flex-col duration-200"
        >
          <TerminalClient
            webSocketUrl={`wss://${window.location.host}/ws/clusters/${clusterId}`}
          />
        </div>
      ) : (
        <p
          key="terminal-unavailable"
          className="fade-in-0 zoom-in-95 animate-in p-4 text-mauve-11 duration-200"
        >
          Terminal is available once the cluster is running.
        </p>
      )}
    </div>
  );
}

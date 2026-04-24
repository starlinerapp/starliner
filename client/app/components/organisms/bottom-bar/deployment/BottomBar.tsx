import React, { memo, useState } from "react";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/client/generated";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";
import Logs from "~/components/organisms/bottom-bar/deployment/Logs";
import TerminalClient from "~/components/atoms/terminal/Terminal.client";

type Deployment =
  | ResponseGitDeployment
  | ResponseImageDeployment
  | ResponseIngressDeployment
  | ResponseDatabaseDeployment;

interface BottomBarProps {
  deployment: Deployment | undefined;
}

const navigationItems = ["Logs", "Terminal"] as const;
type NavigationItem = (typeof navigationItems)[number];

function BottomBarComponent({ deployment }: BottomBarProps) {
  const [selected, setSelected] = useState<NavigationItem>("Logs");

  return (
    <div className="-mt-1 flex h-full flex-col">
      <NavigationBar
        items={navigationItems}
        selected={selected}
        onSelect={setSelected}
      />
      {selected === "Logs" ? (
        <div className="flex min-h-0 flex-1 flex-col p-4">
          <Logs deployment={deployment} />
        </div>
      ) : deployment ? (
        <TerminalClient
          webSocketUrl={`wss://${window.location.host}/ws/deployments/${deployment?.id}`}
        />
      ) : (
        <p className="text-mauve-11 p-4">
          No deployment selected. Select one to connect to the terminal.
        </p>
      )}
    </div>
  );
}

const BottomBar = memo(
  BottomBarComponent,
  (prevProps, nextProps) =>
    prevProps.deployment?.id === nextProps.deployment?.id,
);

export default BottomBar;

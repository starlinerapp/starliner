import { memo, useMemo, useState } from "react";
import TerminalClient from "~/components/atoms/terminal/Terminal.client";
import Logs from "~/components/organisms/bottom-bar/deployment/Logs";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/clients/server/generated";

type Deployment =
  | ResponseGitDeployment
  | ResponseImageDeployment
  | ResponseIngressDeployment
  | ResponseDatabaseDeployment;

interface BottomBarProps {
  deployment: Deployment | undefined;
  showTerminal: boolean;
}

const navigationItems = ["Logs", "Terminal"] as const;
type NavigationItem = (typeof navigationItems)[number];

function BottomBarComponent({ deployment, showTerminal }: BottomBarProps) {
  const [selected, setSelected] = useState<NavigationItem>("Logs");

  const visibleNavigationItems = useMemo(
    () => (showTerminal ? navigationItems : (["Logs"] as const)),
    [showTerminal],
  );

  const activeTab: NavigationItem = showTerminal ? selected : "Logs";

  return (
    <div className="-mt-1 flex h-full flex-col">
      <NavigationBar
        items={visibleNavigationItems}
        selected={activeTab}
        onSelect={setSelected}
      />
      {activeTab === "Logs" ? (
        <div key="logs" className="flex min-h-0 flex-1 flex-col p-4">
          <Logs deployment={deployment} />
        </div>
      ) : deployment && showTerminal ? (
        <div key="terminal" className="flex min-h-0 flex-1 flex-col">
          <TerminalClient
            webSocketUrl={`wss://${window.location.host}/ws/deployments/${deployment.id}`}
          />
        </div>
      ) : (
        <p key="terminal-unavailable" className="p-4 text-mauve-11">
          No deployment selected. Select one to connect to the terminal.
        </p>
      )}
    </div>
  );
}

const BottomBar = memo(
  BottomBarComponent,
  (prevProps, nextProps) =>
    prevProps.deployment?.id === nextProps.deployment?.id &&
    prevProps.showTerminal === nextProps.showTerminal,
);

export default BottomBar;

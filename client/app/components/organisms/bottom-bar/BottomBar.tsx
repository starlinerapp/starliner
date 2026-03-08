import React, { useState } from "react";
import type {
  ResponseDatabaseDeployment,
  ResponseGitDeployment,
  ResponseImageDeployment,
  ResponseIngressDeployment,
} from "~/server/api/client/generated";
import NavigationBar from "../navigation-bar/NavigationBar";
import Logs from "./Logs";

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

export default function BottomBar({ deployment }: BottomBarProps) {
  const [selected, setSelected] = useState<NavigationItem>("Logs");

  return (
    <div className="-mt-1 flex h-full flex-col">
      <NavigationBar
        items={navigationItems}
        selected={selected}
        onSelect={setSelected}
      />
      <div className="min-h-0 flex-1 overflow-y-scroll p-4">
        {selected === "Logs" ? <Logs deployment={deployment} /> : <></>}
      </div>
    </div>
  );
}

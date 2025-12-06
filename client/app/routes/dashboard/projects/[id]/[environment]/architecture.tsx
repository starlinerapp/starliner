import React, { useState } from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";
import ArchitectureCanvas from "~/components/organisms/canvas/ArchitectureCanvas";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";

const services = ["GitHub"] as const;
export type Service = (typeof services)[number];

export default function ProjectArchitecture() {
  const [selected, setSelected] = useState<Service>("GitHub");

  const renderService = () => {
    switch (selected) {
      case "GitHub":
        return <></>;
    }
  };

  return (
    <ResizablePanelGroup direction="horizontal" className="h-full">
      <ResizablePanel
        defaultSize={70}
        className="border-mauve-6 h-full border-r-1"
      >
        <ArchitectureCanvas />
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={30} className="flex h-full flex-col">
        <NavigationBar<Service>
          items={services}
          selected={selected}
          onSelect={setSelected}
        />
        <div className="p-4">{renderService()}</div>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

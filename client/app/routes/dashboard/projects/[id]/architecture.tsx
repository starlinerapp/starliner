import React from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "~/components/atoms/resizable/Resizable";

export default function ProjectArchitecture() {
  return (
    <ResizablePanelGroup direction="horizontal" className="h-full">
      <ResizablePanel
        defaultSize={70}
        className="border-mauve-6 h-full border-r-1 py-3"
      ></ResizablePanel>
      <ResizableHandle />
      <ResizablePanel defaultSize={30} className="h-full"></ResizablePanel>
    </ResizablePanelGroup>
  );
}

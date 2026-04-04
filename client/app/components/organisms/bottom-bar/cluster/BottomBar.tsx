import React, { useState } from "react";
import NavigationBar from "~/components/organisms/navigation-bar/NavigationBar";

const navigationItems = ["Terminal"] as const;
type NavigationItem = (typeof navigationItems)[number];

export default function BottomBar() {
  const [selected, setSelected] = useState<NavigationItem>("Terminal");

  return (
    <div className="-mt-1 flex h-full flex-col">
      <NavigationBar
        items={navigationItems}
        selected={selected}
        onSelect={setSelected}
      />
    </div>
  );
}

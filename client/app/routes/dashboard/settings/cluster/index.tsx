import React from "react";
import HetznerCredentials from "~/components/organisms/settings/HetznerCredentials";

export default function ClusterSettings() {
  return (
    <div className="flex flex-col px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="pt-1 text-xl font-bold">API Keys</h1>
      </div>
      <div className="flex flex-col gap-12 pt-12">
        <HetznerCredentials />
      </div>
    </div>
  );
}


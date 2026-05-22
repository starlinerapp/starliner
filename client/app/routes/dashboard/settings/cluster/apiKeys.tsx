import React from "react";
import HetznerCredentials from "../../../../components/organisms/settings/cluster/HetznerCredentials";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";

export default function ClusterApiKeysSettings() {
  return (
    <>
      <Breadcrumbs
        crumbs={[
          { label: "Settings" },
          { label: "Cluster" },
          { label: "API Keys" },
        ]}
      />
      <div className="flex flex-col px-4 py-4">
        <HetznerCredentials />
      </div>
    </>
  );
}

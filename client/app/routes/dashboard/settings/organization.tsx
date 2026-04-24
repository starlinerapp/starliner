import React from "react";
import OrganizationGeneral from "~/components/organisms/settings/OrganizationGeneral";

export default function OrganizationSettings() {
  return (
    <div className="flex flex-col px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="pt-1 text-xl font-bold">General</h1>
      </div>
      <div className="flex flex-col gap-4 pt-12">
        <OrganizationGeneral />
      </div>
    </div>
  );
}

import React from "react";
import OrganizationTeams from "~/components/organisms/settings/OrganizationTeams";

export default function Teams() {
  return (
    <div className="flex flex-col px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="pt-1 text-xl font-bold">Teams</h1>
      </div>
      <div className="flex flex-col gap-4 pt-10">
        <OrganizationTeams />
      </div>
    </div>
  );
}

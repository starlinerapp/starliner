import React from "react";
import OrganizationGeneral from "~/components/organisms/settings/OrganizationGeneral";
import OrganizationTeams from "~/components/organisms/settings/OrganizationTeams";
import OrganizationInvite from "~/components/organisms/settings/OrganizationInvite";

export default function OrganizationSettings() {
  return (
    <div className="flex flex-col px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="pt-1 text-xl font-bold">Organization Settings</h1>
      </div>
      <div className="flex flex-col gap-4 pt-10">
        <OrganizationGeneral />
        <OrganizationInvite />
        <OrganizationTeams />
      </div>
    </div>
  );
}

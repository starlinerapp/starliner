import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import OrganizationTeams from "~/components/organisms/settings/organization/OrganizationTeams";

export default function Teams() {
  return (
    <>
      <Breadcrumbs
        crumbs={[
          { label: "Settings" },
          { label: "Organization" },
          { label: "Teams" },
        ]}
      />
      <div className="flex w-full flex-col p-4">
        <OrganizationTeams />
      </div>
    </>
  );
}

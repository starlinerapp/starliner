import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import ManageRunners from "~/components/organisms/settings/cluster/ManageRunners";

export default function Runners() {
  return (
    <>
      <Breadcrumbs
        crumbs={[
          { label: "Settings" },
          { label: "Cluster" },
          { label: "Runners" },
        ]}
      />
      <div className="flex w-full flex-col p-4">
        <ManageRunners />
      </div>
    </>
  );
}

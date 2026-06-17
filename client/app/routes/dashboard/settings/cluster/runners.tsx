import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";

export default function Runners() {
  return (
    <Breadcrumbs
      crumbs={[
        { label: "Settings" },
        { label: "Cluster" },
        { label: "Runners" },
      ]}
    />
  );
}

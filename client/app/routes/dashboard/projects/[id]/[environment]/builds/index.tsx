import React from "react";
import LogsCard from "~/components/organisms/logs-card/LogsCard";

export default function Builds() {
  return (
    <div className="flex flex-col gap-4 p-4">
      <LogsCard serviceName="Frontend" />
      <LogsCard serviceName="Backend" />
      <LogsCard serviceName="Frontend" />
      <LogsCard serviceName="Frontend" />
    </div>
  );
}

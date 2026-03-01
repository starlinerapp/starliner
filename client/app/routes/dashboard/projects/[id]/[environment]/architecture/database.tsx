import React from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight } from "~/components/atoms/icons";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";

export default function Database() {
  const trpc = useTRPC();

  const createDatabaseMutation = useMutation(
    trpc.deployment.deployDatabase.mutationOptions(),
  );

  const { environment: currentEnvironment } = useEnvironment();

  function handleDeployClicked() {
    if (!currentEnvironment) return;

    createDatabaseMutation.mutate({
      id: currentEnvironment.id,
    });
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <p>PostgreSQL</p>
        <p className="text-mauve-11 truncate text-sm">
          Powerful, open source relational database
        </p>
      </div>
      <Button
        size="sm"
        className="w-28 flex-shrink-0 py-1.5"
        onClick={handleDeployClicked}
      >
        Deploy
        <ArrowRight className="w-4 stroke-2" />
      </Button>
    </div>
  );
}

import React, { useMemo } from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, Postgres } from "~/components/atoms/icons";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

export default function Database() {
  const { id, environment } = useParams();
  const trpc = useTRPC();

  const { data: project } = useQuery(
    trpc.project.getProject.queryOptions({ id: Number(id) }),
  );

  const currentEnvironment = useMemo(
    () => project?.environments.find((e) => e.slug === environment),
    [project, environment],
  );

  const createDatabaseMutation = useMutation(
    trpc.environment.deployDatabase.mutationOptions(),
  );

  function handleDeployClicked() {
    if (!currentEnvironment) return;

    createDatabaseMutation.mutate({
      id: currentEnvironment.id,
    });
  }

  return (
    <div className="border-mauve-6 flex max-w-full min-w-[350px] items-center justify-between gap-4 overflow-hidden rounded-md border px-4 py-3 text-sm">
      <div className="flex min-w-0 flex-1 items-center gap-4">
        <Postgres className="h-9 w-9 flex-shrink-0" />
        <div className="flex min-w-0 flex-col gap-0.5">
          <p className="truncate font-medium">PostgreSQL</p>
          <p className="text-mauve-11 truncate text-xs">
            Powerful, open source relational database
          </p>
        </div>
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

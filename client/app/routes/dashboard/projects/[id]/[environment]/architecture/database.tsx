import React from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight, Postgres } from "~/components/atoms/icons";

export default function Database() {
  function handleDeployClicked() {}

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

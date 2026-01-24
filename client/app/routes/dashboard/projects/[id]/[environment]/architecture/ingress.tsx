import React from "react";
import { ArrowRight, Traefik } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";

export default function Ingress() {
  return (
    <div className="border-mauve-6 flex max-w-full min-w-[350px] items-center justify-between gap-4 overflow-hidden rounded-md border px-4 py-3 text-sm">
      <div className="flex min-w-0 flex-1 items-center gap-4">
        <Traefik className="h-9 w-9 flex-shrink-0" />
        <div className="flex min-w-0 flex-col gap-0.5">
          <p className="truncate font-medium">Traefik</p>
          <p className="text-mauve-11 truncate text-xs">
            Make your HTTP(S) network service available
          </p>
        </div>
      </div>

      <Button
        size="sm"
        className="w-28 flex-shrink-0 py-1.5"
        onClick={() => {
          console.log("Deploy clicked");
        }}
      >
        Deploy
        <ArrowRight className="w-4 stroke-2" />
      </Button>
    </div>
  );
}

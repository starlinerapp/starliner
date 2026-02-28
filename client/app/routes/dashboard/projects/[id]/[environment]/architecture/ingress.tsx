import React, { useState } from "react";
import { ArrowRight, Plus } from "~/components/atoms/icons";
import Button from "~/components/atoms/button/Button";
import { useTRPC } from "~/utils/trpc/react";
import { useMutation } from "@tanstack/react-query";
import { useEnvironment } from "~/routes/dashboard/projects/[id]/[environment]/architecture/layout";

interface Path {
  path: string;
  pathType: string;
  service: string;
}

interface Host {
  name: string;
  paths: Path[];
}

export default function Ingress() {
  const trpc = useTRPC();

  const createIngressMutation = useMutation(
    trpc.deployment.deployIngress.mutationOptions(),
  );

  const { environment: currentEnvironment } = useEnvironment();

  function handleDeployClicked() {
    if (!currentEnvironment) return;
    createIngressMutation.mutate({
      id: currentEnvironment.id,
    });
  }

  const emptyPathEntry = {
    path: "",
    pathType: "",
    service: "",
  };

  const emptyHostEntry = {
    name: "",
    paths: [emptyPathEntry],
  };

  const [hosts, setHosts] = useState<Host[]>([emptyHostEntry]);

  function handleAddPathClicked(hostIndex: number) {
    setHosts((prev) =>
      prev.map((h, i) =>
        i === hostIndex ? { ...h, paths: [...h.paths, emptyPathEntry] } : h,
      ),
    );
  }

  function handleAddHostClicked() {
    setHosts((prev) => [...prev, emptyHostEntry]);
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-1">
        <p>Traefik</p>
        <p className="text-mauve-11 truncate text-sm">
          Make your HTTP(S) network service available
        </p>
      </div>
      <div className="flex flex-col gap-1">
        <p className="text-sm">Hosts</p>
        <div className="flex flex-col gap-3">
          {hosts.map((host, hostIndex) => (
            <div key={hostIndex} className="flex flex-col gap-1">
              <div className="flex gap-2">
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
                  type="text"
                  value={host.name}
                  placeholder="Host*"
                />
              </div>
              <div className="border-mauve-6 relative flex flex-col gap-3 border-l-2 pl-6">
                {host.paths.map((path, pathIndex) => (
                  <div key={pathIndex} className="relative flex flex-col gap-1">
                    {/* L connector */}
                    <div className="border-mauve-6 absolute -left-6.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />

                    <div className="flex w-full gap-1">
                      <input
                        className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-32 rounded-md border p-2 text-sm"
                        type="text"
                        placeholder="Path Type*"
                      />
                      <input
                        className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-24 rounded-md border p-2 text-sm"
                        type="text"
                        placeholder="Path*"
                      />
                    </div>

                    <div className="border-mauve-6 relative ml-1 flex flex-col gap-1 border-l-2 pl-6">
                      {/* Inner L connector */}
                      <div className="border-mauve-6 absolute -left-0.5 h-6 w-6 rounded-bl-md border-b-2 border-l-2" />

                      <input
                        className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-32 rounded-md border p-2 text-sm"
                        type="text"
                        placeholder="Service*"
                      />
                    </div>
                  </div>
                ))}
              </div>
              <Button
                intent="text"
                className="pl-10"
                onClick={() => handleAddPathClicked(hostIndex)}
              >
                <Plus className="w-3 stroke-3" /> Add Path
              </Button>
            </div>
          ))}
        </div>
        <Button intent="text" onClick={handleAddHostClicked}>
          <Plus className="w-3 stroke-3" /> Add Host
        </Button>
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

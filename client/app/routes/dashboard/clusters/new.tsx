import Button from "~/components/atoms/button/Button";
import React, { useEffect, useState } from "react";
import { type SubmitHandler, useForm } from "react-hook-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTRPC } from "~/utils/trpc/react";
import { useNavigate } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import WarningBanner from "~/components/atoms/banner/WarningBanner";
import { ChevronDown, Cross, Hetzner, K8S } from "~/components/atoms/icons";
import { AnimatePresence, motion } from "framer-motion";
import { cn } from "~/utils/cn";

interface NewClusterFormInput {
  name: string;
  serverType: string;
  teamId: string;
}

export default function NewCluster() {
  const queryClient = useQueryClient();
  const trpc = useTRPC();
  const navigate = useNavigate();

  const organization = useOrganizationContext();

  const [isHetznerOpen, setIsHetznerOpen] = useState(false);

  const { data: hetznerCredentialData, isLoading: isCredentialLoading } =
    useQuery(
      trpc.organization.getHetznerCredential.queryOptions({
        id: organization.id,
      }),
    );

  const isCredentialValid = !!hetznerCredentialData?.credential?.secret;

  const { data: teamsData } = useQuery(
    trpc.team.getUserTeams.queryOptions(
      { organizationId: organization.id },
      { enabled: organization.isOwner },
    ),
  );

  const createClusterMutation = useMutation(
    trpc.cluster.createCluster.mutationOptions({
      onSuccess: async (newCluster) => {
        await queryClient.invalidateQueries({
          queryKey: trpc.organization.getOrganizationClusters.queryKey({
            id: organization.id,
          }),
        });
        navigate(`/${organization.slug}/clusters/${newCluster.id}`);
      },
    }),
  );

  const { register, handleSubmit, watch, setValue } =
    useForm<NewClusterFormInput>();

  const nameInput = watch("name", "");
  const teamIdInput = watch("teamId", "");

  useEffect(() => {
    if (!teamsData?.length) return;

    const defaultTeam =
      teamsData.find((t) => t.slug === organization.slug) ?? teamsData[0];

    setValue("teamId", String(defaultTeam.id));
  }, [teamsData, organization.slug, setValue]);

  const onSubmit: SubmitHandler<NewClusterFormInput> = (data) => {
    createClusterMutation.mutate({
      organizationId: organization.id,
      name: data.name,
      serverType: data.serverType,
      teamId: Number(data.teamId),
    });
  };

  return (
    <div className="flex flex-col gap-16 px-8 py-4">
      <div className="flex flex-col gap-8">
        <h1 className="text-xl font-bold">New Cluster</h1>

        <div className="text-mauve-12 flex flex-col gap-4 text-sm">
          <div className="flex flex-col gap-1">
            <p className="text-sm font-bold">Starliner on your local machine</p>

            <p className="text-mauve-11 text-sm">
              Quickly test and validate the Starliner solution on your computer.
            </p>
          </div>

          <div className="border-mauve-6 hover:bg-gray-2 flex h-32 w-96 cursor-pointer flex-col justify-between rounded-md border p-3 shadow-xs transition-colors">
            <span className="flex items-start justify-between">
              <K8S className="h-10 w-10" />

              <p className="bg-violet-10 border-violet-6 rounded-md px-2 py-1 text-xs text-white">
                3 min to setup
              </p>
            </span>

            <p className="font-bold">Local machine (demo)</p>
          </div>
        </div>

        <div className="text-mauve-12 flex flex-col gap-4 text-sm">
          <div className="flex flex-col gap-1">
            <p className="text-sm font-bold">Or choose your hosting mode</p>

            <p className="text-mauve-11 text-sm">
              Manage your infrastructure across different providers.
            </p>
          </div>

          <div
            onClick={() => setIsHetznerOpen((open) => !open)}
            className="border-mauve-6 hover:bg-gray-2 relative flex h-32 w-96 cursor-pointer flex-col rounded-md border p-3 shadow-xs transition-colors"
          >
            <Hetzner className="h-10 w-10" />
            <div className="absolute top-1/2 right-3 -translate-y-1/2">
              <ChevronDown
                className={cn("h-5 w-5", isHetznerOpen && "rotate-180")}
              />
            </div>
            <p className="mt-auto font-bold">Hetzner Cloud</p>
          </div>

          <AnimatePresence initial={false}>
            {isHetznerOpen && (
              <motion.div
                initial={{ opacity: 0, y: -6 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -6 }}
                transition={{ duration: 0.1, ease: "easeOut" }}
                className="overflow-hidden"
              >
                <div className="border-mauve-6 bg-mauve-2 flex flex-col gap-6 rounded-md border p-4 shadow-xs">
                  <div className="flex flex-col gap-3">
                    <div className="flex items-start justify-between">
                      <span className="flex items-center gap-4">
                        <Hetzner className="h-8 w-8" />

                        <p className="font-bold">Hetzner Cloud</p>
                      </span>

                      <span
                        onClick={(e) => {
                          e.stopPropagation();
                          setIsHetznerOpen(false);
                        }}
                        className="hover:bg-gray-3 cursor-pointer rounded-md p-1"
                      >
                        <Cross className="stroke-mauve-11 h-5 w-5" />
                      </span>
                    </div>

                    {isCredentialLoading ? null : isCredentialValid ? null : (
                      <WarningBanner
                        text="You must enter your Hetzner API Key to create a cluster."
                        linkOut={{
                          text: "Organization Settings",
                          href: `/${organization.slug}/settings/organization`,
                        }}
                        className="my-2"
                      />
                    )}
                  </div>

                  <div className="flex flex-col gap-2">
                    <form
                      className="flex flex-col gap-4"
                      onSubmit={handleSubmit(onSubmit)}
                    >
                      <div className="flex items-end gap-2">
                        <div className="flex flex-col gap-1">
                          <label htmlFor="name">Cluster Name*</label>

                          <input
                            id="name"
                            className="border-mauve-6 h-10 w-80 rounded-md border bg-white px-2 py-1 text-sm"
                            type="text"
                            placeholder="Name*"
                            {...register("name")}
                          />
                        </div>

                        <div className="flex flex-col gap-1">
                          <label htmlFor="teamId">Team*</label>

                          <div className="relative h-10 w-52">
                            <select
                              id="teamId"
                              {...register("teamId", { required: true })}
                              className="border-mauve-6 h-full w-full appearance-none rounded-md border bg-white px-2 py-1 pr-8 text-sm"
                              disabled={!teamsData?.length}
                            >
                              <option value="" disabled>
                                Team*
                              </option>

                              {teamsData?.map((team) => (
                                <option key={team.id} value={team.id}>
                                  {team.slug}
                                </option>
                              ))}
                            </select>

                            <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                              <ChevronDown
                                width={15}
                                className="stroke-mauve-10"
                              />
                            </div>
                          </div>
                        </div>
                      </div>

                      <div className="flex items-end gap-2">
                        <div className="flex flex-col gap-1">
                          <label htmlFor="serverType">Server Type*</label>

                          <div className="relative h-10 w-52">
                            <select
                              id="serverType"
                              className="border-mauve-6 h-full w-full appearance-none rounded-md border bg-white px-2 py-1 pr-8 text-sm"
                              defaultValue="cx23"
                              {...register("serverType", {
                                required: true,
                              })}
                            >
                              <option value="cx23">CX23</option>
                              <option value="cpx22">CPX22</option>
                            </select>

                            <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
                              <ChevronDown
                                width={15}
                                className="stroke-mauve-10"
                              />
                            </div>
                          </div>
                        </div>
                      </div>

                      <Button
                        className="h-10 w-32"
                        disabled={!nameInput || !teamIdInput}
                        type="submit"
                      >
                        Create Cluster
                      </Button>
                    </form>
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
}

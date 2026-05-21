import Button from "~/components/atoms/button/Button";
import React, { useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { CardSkeleton } from "~/components/atoms/card/CardSkeleton";
import { Card } from "~/components/atoms/card/Card";
import {
  ArrowRight,
  FolderOpen,
  MagnifyingGlass,
} from "~/components/atoms/icons";
import { motion } from "framer-motion";
import { formatDistanceToNow } from "date-fns";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import { PlusIcon } from "lucide-react";

export default function Clusters() {
  const trpc = useTRPC();

  const navigate = useNavigate();
  const { slug } = useParams();

  const [search, setSearch] = useState("");

  const organization = useOrganizationContext();

  const { data: clustersData, isLoading } = useQuery(
    trpc.organization.getOrganizationClusters.queryOptions({
      id: organization.id,
    }),
  );

  const filteredClusters = useMemo(() => {
    const query = search.trim().toLowerCase();

    if (!query) return clustersData ?? [];

    return (clustersData ?? []).filter((cluster) => {
      return (
        cluster.name.toLowerCase().includes(query) ||
        cluster.teamSlugs.some((teamSlug) =>
          teamSlug.toLowerCase().includes(query),
        )
      );
    });
  }, [clustersData, search]);

  const arrowVariants = {
    rest: { x: 0 },
    hover: { x: 2 },
  };

  return (
    <>
      <Breadcrumbs crumbs={[{ label: "All Clusters" }]} />
      <div className="flex flex-col gap-8 p-4">
        <div className="flex flex-col gap-3">
          <div className="flex gap-2">
            <div className="relative flex-1">
              <MagnifyingGlass className="text-mauve-11 absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 stroke-2" />
              <input
                type="text"
                value={search}
                onChange={(event) => setSearch(event.target.value)}
                className="border-mauve-6 bg-gray-3 w-full rounded-md border py-2 pr-2 pl-9 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                placeholder="Search for clusters"
              />
            </div>

            <Button
              className="flex w-36 items-center gap-1"
              onClick={() => navigate(`/${slug}/clusters/new`)}
            >
              <PlusIcon className="h-4 w-4" />
              Create Cluster
            </Button>
          </div>
          {isLoading ? (
            <div className="grid grid-cols-[repeat(auto-fit,minmax(350px,1fr))] gap-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <CardSkeleton key={i} />
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-[repeat(auto-fit,minmax(350px,1fr))] justify-start gap-4">
              {filteredClusters.map((cluster, i) => (
                <Link to={`/${slug}/clusters/${cluster.id}`} key={i}>
                  <motion.div initial="rest" animate="rest" whileHover="hover">
                    <Card>
                      <div className="flex h-full flex-col">
                        <div className="flex items-center rounded-t-md px-4 pt-2 pb-1">
                          <div className="flex w-full justify-between">
                            <div>
                              <FolderOpen className="fill-mauve-11 w-6" />{" "}
                              <h2 className="text-mauve-12 font-semibold">
                                {cluster.name}
                              </h2>
                            </div>
                            <motion.div
                              variants={arrowVariants}
                              transition={{
                                type: "spring",
                                stiffness: 500,
                                damping: 30,
                              }}
                            >
                              <ArrowRight className="text-mauve-11 w-5 pt-2" />
                            </motion.div>
                          </div>
                        </div>

                        <div className="flex h-full flex-col gap-2 px-4 pb-4">
                          <p className="text-mauve-11 text-xs">
                            Created{" "}
                            <span>
                              {formatDistanceToNow(
                                new Date(cluster.createdAt),
                                {
                                  addSuffix: true,
                                },
                              )}
                            </span>
                          </p>
                          <div className="flex flex-wrap gap-1.5">
                            {cluster.teamSlugs.map((teamSlug, i) => (
                              <p
                                key={i}
                                className="text-violet-11 bg-violet-3 w-fit rounded-md px-2 py-1 text-xs"
                              >
                                #<span>{teamSlug}</span>
                              </p>
                            ))}
                          </div>
                        </div>
                      </div>
                    </Card>
                  </motion.div>
                </Link>
              ))}
              <div></div>
              <div></div>
              <div></div>
              <div></div>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

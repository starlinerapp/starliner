import { useQuery } from "@tanstack/react-query";
import { formatDistanceToNow } from "date-fns";
import { motion } from "framer-motion";
import { PlusIcon } from "lucide-react";
import { useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import Button from "~/components/atoms/button/Button";
import { Card } from "~/components/atoms/card/Card";
import { CardSkeleton } from "~/components/atoms/card/CardSkeleton";
import {
  ArrowRight,
  FolderOpen,
  MagnifyingGlass,
} from "~/components/atoms/icons";
import Breadcrumbs from "~/components/organisms/breadcrumbs/Breadcrumbs";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";

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
              <MagnifyingGlass className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 stroke-2 text-mauve-11" />
              <input
                type="text"
                value={search}
                onChange={(event) => setSearch(event.target.value)}
                className="w-full rounded-md border border-mauve-6 bg-gray-2 py-2 pr-2 pl-9 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)] placeholder:text-mauve-11"
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
                              <FolderOpen className="w-6 fill-mauve-11" />{" "}
                              <h2 className="font-semibold text-mauve-12">
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
                              <ArrowRight className="w-5 pt-2 text-mauve-11" />
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
                                className="w-fit rounded-md bg-violet-3 px-2 py-1 text-violet-11 text-xs"
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

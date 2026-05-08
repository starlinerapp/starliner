import Button from "~/components/atoms/button/Button";
import React from "react";
import { Link, useNavigate, useParams } from "react-router";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { CardSkeleton } from "~/components/atoms/card/CardSkeleton";
import { Card } from "~/components/atoms/card/Card";
import { ArrowRight, FolderOpen } from "~/components/atoms/icons";
import { motion } from "framer-motion";
import { formatDistanceToNow } from "date-fns";

export default function Clusters() {
  const trpc = useTRPC();

  const navigate = useNavigate();
  const { slug } = useParams();

  const organization = useOrganizationContext();

  const { data: clustersData, isLoading } = useQuery(
    trpc.organization.getOrganizationClusters.queryOptions({
      id: organization.id,
    }),
  );

  const arrowVariants = {
    rest: { x: 0 },
    hover: { x: 2 },
  };

  return (
    <div className="flex flex-col gap-8 px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="text-xl font-bold">Clusters</h1>
        <Button
          className="w-32"
          onClick={() => navigate(`/${slug}/clusters/new`)}
        >
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
        <div className="grid grid-cols-[repeat(auto-fit,minmax(350px,350px))] justify-start gap-4">
          {clustersData?.map((cluster, i) => (
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
                          {formatDistanceToNow(new Date(cluster.createdAt), {
                            addSuffix: true,
                          })}
                        </span>
                      </p>
                      <p className="text-mauve-11 bg-violet-3 border-mauve-6 w-fit rounded-md border px-2 py-1 text-xs">
                        #<span>{cluster.teamSlug}</span>
                      </p>
                    </div>
                  </div>
                </Card>
              </motion.div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

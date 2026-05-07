import React from "react";
import Button from "~/components/atoms/button/Button";
import { Link, useNavigate, useParams } from "react-router";
import { Card } from "~/components/atoms/card/Card";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { CardSkeleton } from "~/components/atoms/card/CardSkeleton";
import { ArrowRight, FolderOpen } from "~/components/atoms/icons";
import { formatDistanceToNow } from "date-fns";
import { motion } from "framer-motion";

export default function Projects() {
  const navigate = useNavigate();
  const { slug } = useParams();

  const organization = useOrganizationContext();

  const trpc = useTRPC();
  const { data: projectsData, isLoading } = useQuery(
    trpc.organization.getUserProjects.queryOptions({
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
        <h1 className="text-xl font-bold">Projects</h1>
        <Button
          className="w-32"
          onClick={() => navigate(`/${slug}/projects/new`)}
        >
          Create Project
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
          {projectsData?.map((project, i) => (
            <Link to={`/${slug}/projects/${project.id}`} key={i}>
              <motion.div initial="rest" animate="rest" whileHover="hover">
                <Card>
                  <div className="flex h-full flex-col">
                    <div className="flex items-center rounded-t-md px-4 pt-2 pb-1">
                      <div className="flex w-full justify-between">
                        <div>
                          <FolderOpen className="fill-mauve-11 w-6" />{" "}
                          <h2 className="text-mauve-12 font-semibold">
                            {project.name}
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

                    <div className="flex h-full flex-col justify-between px-4 pb-4">
                      <p className="text-mauve-11 text-xs">
                        #<span>{project.teamSlug}</span>
                      </p>
                      <p className="text-mauve-11 text-xs">
                        Created{" "}
                        <span>
                          {formatDistanceToNow(new Date(project.createdAt), {
                            addSuffix: true,
                          })}
                        </span>
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

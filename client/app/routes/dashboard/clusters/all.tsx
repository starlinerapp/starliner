import Button from "~/components/atoms/button/Button";
import React from "react";
import { useNavigate, useParams } from "react-router";

export default function Clusters() {
  const navigate = useNavigate();
  const { slug } = useParams();

  return (
    <div className="px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="text-xl font-bold">Clusters</h1>
        <Button
          className="w-32"
          onClick={() => navigate(`/${slug}/clusters/new`)}
        >
          Create Cluster
        </Button>
      </div>
    </div>
  );
}

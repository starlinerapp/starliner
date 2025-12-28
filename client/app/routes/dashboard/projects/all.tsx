import React from "react";
import Button from "~/components/atoms/button/Button";
import { useNavigate, useParams } from "react-router";

export default function Projects() {
  const navigate = useNavigate();
  const { slug } = useParams();

  return (
    <div className="px-8 py-4">
      <div className="flex w-full items-center justify-between">
        <h1 className="text-xl font-bold">Projects</h1>
        <Button
          className="w-32"
          onClick={() => navigate(`/${slug}/projects/new`)}
        >
          Create Project
        </Button>
      </div>
    </div>
  );
}

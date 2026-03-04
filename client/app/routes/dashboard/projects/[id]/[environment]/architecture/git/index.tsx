import React from "react";
import Button from "~/components/atoms/button/Button";
import { ArrowRight } from "~/components/atoms/icons";

export default function Git() {
  return (
    <>
      <form className="flex flex-col gap-4">
        <div className="flex flex-col gap-1">
          <p>Import Git Repository</p>
          <p className="text-mauve-11 truncate text-sm">
            Enter a Git repository URL to deploy.
          </p>
        </div>
        <div className="flex flex-col gap-2">
          <div className="flex flex-col gap-1">
            <p className="text-sm">URL</p>
            <div className="flex gap-2">
              <input
                className="border-mauve-6 disabled:text-mauve-10 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm disabled:hover:cursor-not-allowed"
                type="text"
                placeholder="Git URL*"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1">
            <div className="flex items-center gap-2">
              <div className="w-full">
                <p className="text-sm">Service Name</p>
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-52 rounded-md border-1 p-2 text-sm"
                  type="text"
                  placeholder="Name*"
                />
              </div>
              <div className="w-full">
                <p className="text-sm">Dockerfile</p>
                <input
                  className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 w-full min-w-24 rounded-md border-1 p-2 text-sm"
                  type="text"
                  placeholder="Path*"
                />
              </div>
            </div>
          </div>
        </div>
        <Button type="submit" size="sm" className="w-28 flex-shrink-0 py-1.5">
          Deploy
          <ArrowRight className="w-4 stroke-2" />
        </Button>
      </form>
    </>
  );
}

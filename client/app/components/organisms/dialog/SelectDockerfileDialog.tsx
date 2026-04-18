import React, { useState } from "react";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { useTRPC } from "~/utils/trpc/react";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import { useQuery } from "@tanstack/react-query";
import { cn } from "~/utils/cn";
import { ChevronRight } from "lucide-react";

interface SelectDockerfileDialogProps {
  repositoryOwner: string;
  repositoryName: string;
  path: string;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: (filePath: string) => void;
}

export default function SelectDockerfileDialog({
  repositoryOwner,
  repositoryName,
  path,
  isOpen,
  onOpenChange,
  onConfirm,
}: SelectDockerfileDialogProps) {
  const organizationContext = useOrganizationContext();
  const trpc = useTRPC();
  const [selectedPath, setSelectedPath] = useState<string | null>(null);

  const { data: repositoryContentData, isLoading: isRootLoading } = useQuery(
    trpc.github.getRepositoryFiles.queryOptions(
      {
        organizationId: organizationContext.id,
        owner: repositoryOwner,
        repo: repositoryName,
        path: path,
      },
      {
        enabled: repositoryOwner !== "" && repositoryName !== "",
      },
    ),
  );

  function handleContinue() {
    if (selectedPath === null) return;
    const normalizedPath = path.startsWith("./") ? path.slice(2) : path;
    const relativePath = selectedPath.startsWith(normalizedPath + "/")
      ? "./" + selectedPath.slice(normalizedPath.length + 1)
      : "./" + selectedPath;

    onConfirm(relativePath);
    onOpenChange(false);
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="flex flex-col justify-between gap-4">
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <h1 className="text-mauve-12">Select Dockerfile</h1>
              <p className="text-mauve-11 text-sm">
                If your repository contains multiple Dockerfiles, select the one
                you&#39;d like to use for this deployment.
              </p>
            </div>
            <form className="bg-mauve-2 border-mauve-6 flex max-h-[500px] flex-col gap-2 overflow-y-auto rounded-md border-1 p-2">
              {isRootLoading ? (
                <span className="flex flex-col gap-2">
                  <DirectoryItemSkeleton depth={0} />
                  <DirectoryItemSkeleton depth={0} />
                  <DirectoryItemSkeleton depth={0} />
                </span>
              ) : (
                repositoryContentData?.map((item, i) => (
                  <span key={i}>
                    <DirectoryItem
                      name={item.name}
                      path={item.path ?? item.name}
                      type={item.type}
                      depth={0}
                      repositoryOwner={repositoryOwner}
                      repositoryName={repositoryName}
                      organizationId={organizationContext.id}
                      selectedPath={selectedPath}
                      onSelect={setSelectedPath}
                    />
                  </span>
                ))
              )}
            </form>
          </div>
          <div className="flex w-full justify-end gap-2">
            <Button
              intent="secondary"
              className="w-24"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button
              className="w-24"
              disabled={selectedPath === null}
              onClick={handleContinue}
            >
              Continue
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}

interface DirectoryItemProps {
  name: string;
  path: string;
  type: string;
  depth: number;
  repositoryOwner: string;
  repositoryName: string;
  organizationId: number;
  selectedPath: string | null;
  onSelect: (path: string) => void;
  className?: string;
}

function DirectoryItem({
  name,
  path,
  type,
  depth,
  repositoryOwner,
  repositoryName,
  organizationId,
  selectedPath,
  onSelect,
  className,
}: DirectoryItemProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const isDir = type === "dir";

  const trpc = useTRPC();

  const { data: subFiles, isLoading } = useQuery({
    ...trpc.github.getRepositoryFiles.queryOptions({
      organizationId: organizationId,
      owner: repositoryOwner,
      repo: repositoryName,
      path,
    }),
    enabled: isExpanded && isDir,
  });

  return (
    <div className={cn("flex flex-col", className)}>
      <fieldset
        className="text-mauve-12 flex items-center gap-2"
        style={{ paddingLeft: `${depth * 1.5}rem` }}
      >
        {isDir && (
          <button
            type="button"
            onClick={() => setIsExpanded((prev) => !prev)}
            className="text-mauve-11 hover:text-mauve-12 transition-transform outline-none"
          >
            <ChevronRight
              className="h-5 w-5 transition-transform duration-150"
              style={{
                transform: isExpanded ? "rotate(90deg)" : "rotate(0deg)",
              }}
            />
          </button>
        )}
        {!isDir && (
          <input
            type="radio"
            className="outline-none"
            name="projectFile"
            id={path}
            value={path}
            checked={selectedPath === path}
            onChange={() => onSelect(path)}
          />
        )}
        <label htmlFor={path} className="text-mauve-11 font-mono text-sm">
          {name}
        </label>
      </fieldset>

      {isDir && isExpanded && (
        <div className="flex flex-col gap-2">
          {isLoading ? (
            <span className="flex flex-col gap-2 pl-0.5">
              <DirectoryItemSkeleton depth={depth + 1} className="pt-2" />
              <DirectoryItemSkeleton depth={depth + 1} />
            </span>
          ) : (
            (subFiles?.length ?? 0) > 0 &&
            subFiles?.map((dir, i) => (
              <DirectoryItem
                className={cn(i === 0 && "pt-2")}
                key={i}
                name={dir.name}
                path={dir.path ?? `${path}/${dir.name}`}
                type={dir.type}
                depth={depth + 1}
                repositoryOwner={repositoryOwner}
                repositoryName={repositoryName}
                organizationId={organizationId}
                selectedPath={selectedPath}
                onSelect={onSelect}
              />
            ))
          )}
        </div>
      )}
    </div>
  );
}

interface DirectoryItemSkeletonProps {
  depth?: number;
  className?: string;
}

function DirectoryItemSkeleton({
  depth = 0,
  className,
}: DirectoryItemSkeletonProps) {
  return (
    <div
      className={cn("flex items-center gap-2", className)}
      style={{ paddingLeft: `${depth * 1.5}rem` }}
    >
      <div className="bg-mauve-5 h-4 w-4 animate-pulse rounded-full" />
      <div
        className="bg-mauve-5 h-4 animate-pulse rounded"
        style={{ width: `${Math.floor(Math.random() * 40) + 40}%` }}
      />
    </div>
  );
}

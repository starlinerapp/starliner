import React, { useState } from "react";
import { Dialog, DialogContent } from "~/components/atoms/dialog/Dialog";
import Button from "~/components/atoms/button/Button";
import { ChevronRight } from "lucide-react";
import { useTRPC } from "~/utils/trpc/react";
import { useQuery } from "@tanstack/react-query";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import { cn } from "~/utils/cn";

interface SelectProjectDirectoryDialogProps {
  repositoryOwner: string;
  repositoryName: string;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: (projectDirectoryPath: string) => void;
}

export default function SelectProjectDirectoryDialog({
  repositoryOwner,
  repositoryName,
  isOpen,
  onOpenChange,
  onConfirm,
}: SelectProjectDirectoryDialogProps) {
  const organizationContext = useOrganizationContext();
  const trpc = useTRPC();
  const [selectedPath, setSelectedPath] = useState<string | null>(null);
  const [isRootExpanded, setIsRootExpanded] = useState(true);

  const { data: repositoryContentData, isLoading: isRootLoading } = useQuery(
    trpc.github.getRepositoryFiles.queryOptions(
      {
        organizationId: organizationContext.id,
        owner: repositoryOwner,
        repo: repositoryName,
        path: "",
      },
      {
        enabled: repositoryOwner !== "" && repositoryName !== "",
      },
    ),
  );

  const projectDirectories =
    repositoryContentData?.filter((file) => file.type === "dir") ?? [];

  function handleContinue() {
    if (selectedPath === null) return;
    onConfirm(selectedPath === "" ? "./" : `./${selectedPath}`);
    onOpenChange(false);
  }

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="flex flex-col justify-between gap-4">
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <h1 className="text-mauve-12">Select Project Directory</h1>
              <p className="text-mauve-11 text-sm">
                Select the directory containing your source code. For monorepos,
                create a separate deployment for each directory you want to
                deploy.
              </p>
            </div>
            <form className="bg-mauve-2 border-mauve-6 flex max-h-125 flex-col gap-2 overflow-y-auto rounded-md border p-2">
              <fieldset className="text-mauve-12 flex items-center gap-2">
                <button
                  type="button"
                  onClick={() => setIsRootExpanded((prev) => !prev)}
                  className="text-mauve-11 hover:text-mauve-12 outline-none"
                >
                  <ChevronRight
                    className="h-5 w-5 transition-transform duration-150"
                    style={{
                      transform: isRootExpanded
                        ? "rotate(90deg)"
                        : "rotate(0deg)",
                    }}
                  />
                </button>
                <input
                  type="radio"
                  className="outline-none"
                  name="projectDirectory"
                  id="root"
                  value=""
                  checked={selectedPath === ""}
                  onChange={() => setSelectedPath("")}
                />
                <label
                  htmlFor="root"
                  className="text-mauve-11 font-mono text-sm"
                >
                  ./
                </label>
              </fieldset>
              {isRootExpanded &&
                (isRootLoading ? (
                  <span className="flex flex-col gap-2">
                    <DirectoryItemSkeleton depth={1} />
                    <DirectoryItemSkeleton depth={1} />
                    <DirectoryItemSkeleton depth={1} />
                  </span>
                ) : (
                  projectDirectories.map((directory, i) => (
                    <span key={i}>
                      <DirectoryItem
                        name={directory.name}
                        path={directory.path ?? directory.name}
                        depth={1}
                        repositoryOwner={repositoryOwner}
                        repositoryName={repositoryName}
                        organizationId={organizationContext.id}
                        selectedPath={selectedPath}
                        onSelect={setSelectedPath}
                      />
                    </span>
                  ))
                ))}
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
  depth,
  repositoryOwner,
  repositoryName,
  organizationId,
  selectedPath,
  onSelect,
  className,
}: DirectoryItemProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const trpc = useTRPC();

  const { data: subFiles, isLoading } = useQuery({
    ...trpc.github.getRepositoryFiles.queryOptions({
      organizationId: organizationId,
      owner: repositoryOwner,
      repo: repositoryName,
      path,
    }),
    enabled: isExpanded,
  });

  const subDirectories = subFiles?.filter((file) => file.type === "dir") ?? [];

  return (
    <div className={cn("flex flex-col", className)}>
      <fieldset
        className="text-mauve-12 flex items-center gap-2"
        style={{ paddingLeft: `${depth * 1.5}rem` }}
      >
        <button
          type="button"
          onClick={() => setIsExpanded((prev) => !prev)}
          className="text-mauve-11 hover:text-mauve-12 transition-transform outline-none"
        >
          <ChevronRight
            className="h-5 w-5 transition-transform duration-150"
            style={{ transform: isExpanded ? "rotate(90deg)" : "rotate(0deg)" }}
          />
        </button>
        <input
          type="radio"
          className="outline-none"
          name="projectDirectory"
          id={path}
          value={path}
          checked={selectedPath === path}
          onChange={() => onSelect(path)}
        />
        <label htmlFor={path} className="font-mono text-sm">
          {name}
        </label>
      </fieldset>

      {isExpanded && (
        <div className="flex flex-col gap-2">
          {isLoading ? (
            <span className="flex flex-col gap-2 pl-0.5">
              <DirectoryItemSkeleton depth={depth + 1} className="pt-2" />
              <DirectoryItemSkeleton depth={depth + 1} />
            </span>
          ) : (
            subDirectories.length > 0 &&
            subDirectories.map((dir, i) => (
              <DirectoryItem
                className={cn(i === 0 && "pt-2")}
                key={i}
                name={dir.name}
                path={dir.path ?? `${path}/${dir.name}`}
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

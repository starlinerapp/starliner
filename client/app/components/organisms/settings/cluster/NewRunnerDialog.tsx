import { useMemo } from "react";
import { ChevronDown } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

const RUNNER_REPO_URL = "https://github.com/starlinerapp/runner";
const RUNNER_VERSION = "v0.0.1";
const RUNNER_PACKAGE = `runner-${RUNNER_VERSION}-linux-amd64.tar`;
const RUNNER_DIR = `runner-${RUNNER_VERSION}-linux-amd64`;
const RUNNER_DOWNLOAD_URL = `${RUNNER_REPO_URL}/releases/download/${RUNNER_VERSION}/${RUNNER_PACKAGE}`;
const RUNNER_TOKEN = "AIOCFPV54MZ3MPTM77HFGCTKFWNJS";

type RunnerImage = "macos" | "linux" | "windows";

const RUNNER_IMAGE_OPTIONS: {
  value: RunnerImage;
  label: string;
  disabled?: boolean;
}[] = [
  { value: "linux", label: "Linux" },
  { value: "macos", label: "macOS", disabled: true },
  { value: "windows", label: "Windows", disabled: true },
];

function getDownloadScript() {
  return `# Download the runner bundle
$ curl -LO ${RUNNER_DOWNLOAD_URL}

$ tar xf ${RUNNER_PACKAGE}
$ cd ${RUNNER_DIR}
$ ./runner install`;
}

function getConfigureScript() {
  return `# Register the runner with your organization
$ ./runner register --token ${RUNNER_TOKEN}

# Start the runner
$ ./runner run`;
}

export default function NewRunnerDialog() {
  const downloadScript = useMemo(() => getDownloadScript(), []);
  const configureScript = useMemo(() => getConfigureScript(), []);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-2">
        <h1>New self-hosted runner</h1>
        <p className="text-mauve-11 text-sm">
          Adding a self-hosted runner requires that you download, configure, and
          run the{" "}
          <a
            href={RUNNER_REPO_URL}
            target="_blank"
            rel="noreferrer"
            className="underline"
          >
            Starliner runner
          </a>
          . The runner connects to Starliner to execute image builds on your own
          infrastructure.
        </p>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <p className="font-semibold text-mauve-12 text-sm">Runner image</p>
          <div className="flex gap-2">
            {RUNNER_IMAGE_OPTIONS.map((option) => (
              <RunnerImageOption
                key={option.value}
                value={option.value}
                label={option.label}
                disabled={option.disabled}
                selected={option.value === "linux"}
              />
            ))}
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <label
            className="font-semibold text-mauve-12 text-sm"
            htmlFor="architecture"
          >
            Architecture
          </label>
          <ArchitectureSelect />
        </div>

        <div className="flex flex-col gap-2">
          <p className="font-semibold text-mauve-12 text-sm">Download</p>
          <ScriptBlock script={downloadScript} />
        </div>

        <div className="flex flex-col gap-2">
          <p className="font-semibold text-mauve-12 text-sm">Configure</p>
          <ScriptBlock script={configureScript} />
        </div>
      </div>
    </div>
  );
}

function RunnerImageOption({
  value,
  label,
  disabled,
  selected,
}: {
  value: RunnerImage;
  label: string;
  disabled?: boolean;
  selected: boolean;
}) {
  return (
    <label
      className={cn(
        "flex min-w-0 flex-1 items-center gap-2 rounded-md border border-mauve-6 px-3 py-2.5",
        selected && "border-violet-9 bg-violet-2",
        disabled ? "cursor-not-allowed opacity-50" : "cursor-default",
      )}
    >
      <input
        type="radio"
        name="runnerImage"
        value={value}
        checked={selected}
        disabled={disabled}
        readOnly
        className="size-4 shrink-0 accent-[#0969da]"
      />
      <span className="truncate text-mauve-12 text-sm">{label}</span>
    </label>
  );
}

function ArchitectureSelect() {
  return (
    <div className="relative w-full">
      <select
        id="architecture"
        disabled
        value="x64"
        className="h-10 w-full cursor-not-allowed appearance-none rounded-md border border-mauve-6 bg-white p-2 text-mauve-12 text-sm"
      >
        <option value="x64">x64</option>
      </select>
      <div className="pointer-events-none absolute inset-y-0 right-2 flex items-center">
        <ChevronDown width={15} className="stroke-mauve-10" />
      </div>
    </div>
  );
}

function ScriptBlock({ script }: { script: string }) {
  return (
    <pre className="w-full whitespace-pre-wrap break-all rounded-md border border-mauve-6 bg-gray-2 p-3 font-mono text-mauve-11 text-xs leading-relaxed">
      {script}
    </pre>
  );
}

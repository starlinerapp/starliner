import React, { useState } from "react";
import { useParams } from "react-router";
import { useOrganizationContext } from "~/contexts/OrganizationContext";
import Button from "~/components/atoms/button/Button";
import CopyToClipboard from "~/components/atoms/copy-to-clipboard/CopyToClipboard";
import { Cross, ExclamationTriangle, Plus } from "~/components/atoms/icons";
import { cn } from "~/utils/cn";

interface CustomDomain {
  id: string;
  domain: string;
  verificationToken: string;
  status: "pending" | "verified";
}

function isApexDomain(domain: string) {
  return domain.split(".").length === 2;
}

function isValidDomain(domain: string) {
  return /^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$/.test(
    domain.trim().toLowerCase(),
  );
}

export default function DomainsPage() {
  const { environment } = useParams<{ environment: string }>();
  const organization = useOrganizationContext();

  const [domains, setDomains] = useState<CustomDomain[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [showInput, setShowInput] = useState(false);
  const [inputError, setInputError] = useState<string | null>(null);

  // TODO: derive from cluster data
  const cnameTarget =
    environment === "production"
      ? `ingress.${organization.slug}.starliner.cloud`
      : `ingress.${organization.slug}.${environment}.starliner.cloud`;

  const handleAdd = () => {
    const domain = inputValue.trim().toLowerCase();
    if (!isValidDomain(domain)) {
      setInputError("Enter a valid domain (e.g. api.mycompany.com)");
      return;
    }
    if (domains.some((d) => d.domain === domain)) {
      setInputError("This domain is already added");
      return;
    }
    setDomains((prev) => [
      ...prev,
      {
        id: Math.random().toString(36).slice(2),
        domain,
        verificationToken: `starliner-${Math.random().toString(36).slice(2, 14)}`,
        status: "pending",
      },
    ]);
    setInputValue("");
    setShowInput(false);
    setInputError(null);
  };

  return (
    <div className="w-full space-y-4 p-4">
      <div className="border-mauve-6 rounded-md border text-sm shadow-xs">
        <div className="border-mauve-6 bg-gray-2 flex h-14 items-center justify-between border-b px-4">
          <span className="text-mauve-12 text-xs uppercase">Custom Domains</span>
          <Button
            intent="secondary"
            size="sm"
            type="button"
            onClick={() => { setShowInput(true); setInputError(null); }}
          >
            <Plus className="w-3 stroke-3" /> Add Domain
          </Button>
        </div>

        {showInput && (
          <div className="border-mauve-6 flex flex-col gap-2 border-b px-4 py-3">
            <div className="flex gap-2">
              <input
                autoFocus
                className="border-mauve-6 placeholder:text-mauve-11 bg-gray-2 flex-1 rounded-md border p-2 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.12)]"
                placeholder="api.mycompany.com"
                value={inputValue}
                onChange={(e) => { setInputValue(e.target.value); setInputError(null); }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") handleAdd();
                  if (e.key === "Escape") { setShowInput(false); setInputValue(""); setInputError(null); }
                }}
              />
              <Button size="sm" type="button" onClick={handleAdd}>Add</Button>
              <Button intent="secondary" size="sm" type="button" onClick={() => { setShowInput(false); setInputValue(""); setInputError(null); }}>
                Cancel
              </Button>
            </div>
            {inputError && <p className="text-red-11 text-xs">{inputError}</p>}
          </div>
        )}

        {domains.length === 0 && !showInput ? (
          <div className="text-mauve-11 flex h-24 items-center justify-center text-sm">
            No custom domains configured for this environment.
          </div>
        ) : (
          domains.map((domain, index) => (
            <div
              key={domain.id}
              className={cn("flex flex-col gap-3 px-4 py-3", index < domains.length - 1 && "border-mauve-6 border-b")}
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex flex-col gap-1">
                  <span className="font-medium">{domain.domain}</span>
                  <span className={cn("flex items-center gap-1 text-xs", domain.status === "verified" ? "text-green-11" : "text-amber-11")}>
                    {domain.status === "pending" && <ExclamationTriangle className="w-3.5" strokeWidth={2} />}
                    {domain.status === "verified" ? "Verified" : "Pending DNS"}
                  </span>
                </div>
                <div className="flex shrink-0 items-center gap-2">
                  {domain.status === "pending" && (
                    <Button
                      intent="secondary"
                      size="sm"
                      type="button"
                      onClick={() =>
                        // TODO: call backend DNS check
                        setDomains((prev) => prev.map((d) => d.id === domain.id ? { ...d, status: "verified" } : d))
                      }
                    >
                      Check DNS
                    </Button>
                  )}
                  <button
                    type="button"
                    onClick={() => setDomains((prev) => prev.filter((d) => d.id !== domain.id))}
                    className="text-mauve-9 hover:text-red-11 cursor-pointer"
                  >
                    <Cross className="w-3.5" strokeWidth={2} />
                  </button>
                </div>
              </div>

              {domain.status === "pending" && (
                <div className="border-mauve-6 flex flex-col gap-3 rounded-md border p-3 text-xs">
                  <p className="text-mauve-11">
                    Add these DNS records at your provider to verify and connect your domain to Starliner.
                  </p>
                  <table className="border-mauve-6 w-full border-collapse overflow-hidden rounded-md border">
                    <thead>
                      <tr className="border-mauve-6 bg-white border-b">
                        {["Type", "Name", "Value"].map((h) => (
                          <th key={h} className="border-mauve-6 text-mauve-11 border-r px-2 py-1.5 text-left font-medium last:border-r-0">
                            {h}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-mauve-6 border-b">
                        <td className="border-mauve-6 border-r px-2 py-1.5 font-mono">TXT</td>
                        <td className="border-mauve-6 border-r"><CopyToClipboard text={`_starliner-verification.${domain.domain}`} className="font-mono" /></td>
                        <td><CopyToClipboard text={domain.verificationToken} className="font-mono" /></td>
                      </tr>
                      <tr>
                        <td className="border-mauve-6 border-r px-2 py-1.5 font-mono">{isApexDomain(domain.domain) ? "A" : "CNAME"}</td>
                        <td className="border-mauve-6 border-r"><CopyToClipboard text={domain.domain} className="font-mono" /></td>
                        <td><CopyToClipboard text={isApexDomain(domain.domain) ? "TODO: cluster IP" : cnameTarget} className="font-mono" /></td>
                      </tr>
                    </tbody>
                  </table>
                  <p className="text-mauve-11">DNS changes may take a few minutes to propagate.</p>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}

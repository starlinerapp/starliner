import { describe, expect, it } from "vitest";
import {
  buildFullIngressHost,
  getIngressHostSuffix,
  isValidIngressHostPrefix,
  parseIngressHostPrefix,
} from "./ingress-host";

describe("ingress-host", () => {
  const orgSlug = "acme";
  const deploymentDomain = "starliner.cloud";

  it("builds suffixes from deployment environment", () => {
    expect(getIngressHostSuffix(orgSlug, "production", deploymentDomain)).toBe(
      ".acme.starliner.cloud",
    );
    expect(getIngressHostSuffix(orgSlug, "staging", deploymentDomain)).toBe(
      ".acme.staging.starliner.cloud",
    );
    expect(getIngressHostSuffix(orgSlug, "local", deploymentDomain)).toBe(
      ".acme.dev.starliner.cloud",
    );
  });

  it("builds and parses full hostnames", () => {
    const full = buildFullIngressHost(
      "api",
      orgSlug,
      "staging",
      deploymentDomain,
    );
    expect(full).toBe("api.acme.staging.starliner.cloud");
    expect(
      parseIngressHostPrefix(full, orgSlug, "staging", deploymentDomain),
    ).toBe("api");
  });

  it("validates host prefixes", () => {
    expect(isValidIngressHostPrefix("api")).toBe(true);
    expect(isValidIngressHostPrefix("my-service")).toBe(true);
    expect(isValidIngressHostPrefix(" API ")).toBe(true);
    expect(isValidIngressHostPrefix("---")).toBe(false);
    expect(isValidIngressHostPrefix("")).toBe(false);
  });
});

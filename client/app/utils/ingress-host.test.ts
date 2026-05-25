import { describe, expect, it } from "vitest";
import {
  buildFullIngressHost,
  getIngressHostDomain,
  getIngressHostSuffix,
  isValidIngressHost,
  isValidIngressHostPrefix,
  parseIngressHostPrefix,
  sanitizeIngressHostPrefix,
} from "./ingress-host";

describe("ingress-host", () => {
  const orgSlug = "acme";

  it("resolves domain tiers", () => {
    expect(
      getIngressHostDomain({
        deploymentEnvironment: "local",
        environmentSlug: "production",
      }),
    ).toBe("local");
    expect(
      getIngressHostDomain({
        deploymentEnvironment: "production",
        environmentSlug: "production",
      }),
    ).toBe("production");
    expect(
      getIngressHostDomain({
        deploymentEnvironment: "staging",
        environmentSlug: "preview",
      }),
    ).toBe("staging");
  });

  it("builds suffixes for each domain tier", () => {
    expect(getIngressHostSuffix(orgSlug, "production")).toBe(
      ".acme.starliner.cloud",
    );
    expect(getIngressHostSuffix(orgSlug, "staging")).toBe(
      ".acme.staging.starliner.cloud",
    );
    expect(getIngressHostSuffix(orgSlug, "local")).toBe(
      ".acme.dev.starliner.cloud",
    );
  });

  it("builds and parses full hostnames", () => {
    const full = buildFullIngressHost("api", orgSlug, "production");
    expect(full).toBe("api.acme.starliner.cloud");
    expect(parseIngressHostPrefix(full, orgSlug)).toBe("api");
  });

  it("validates host prefixes", () => {
    expect(isValidIngressHostPrefix("api")).toBe(true);
    expect(isValidIngressHostPrefix("my-service")).toBe(true);
    expect(isValidIngressHostPrefix("-bad-")).toBe(true);
    expect(isValidIngressHostPrefix("---")).toBe(false);
    expect(isValidIngressHostPrefix("")).toBe(false);
    expect(sanitizeIngressHostPrefix(" API ")).toBe("api");
  });

  it("validates full hostnames against the required suffix", () => {
    expect(
      isValidIngressHost(
        "api.acme.staging.starliner.cloud",
        orgSlug,
        "staging",
      ),
    ).toBe(true);
    expect(
      isValidIngressHost("api.acme.starliner.cloud", orgSlug, "staging"),
    ).toBe(false);
  });
});

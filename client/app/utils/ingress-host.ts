export type IngressHostDomain = "production" | "staging" | "local";

export function getIngressHostDomain(options: {
  isLocal: boolean;
  environmentSlug: string;
}): IngressHostDomain {
  if (options.isLocal) {
    return "local";
  }

  if (options.environmentSlug === "production") {
    return "production";
  }

  return "staging";
}

export function getIngressHostSuffix(
  organizationSlug: string,
  domain: IngressHostDomain,
): string {
  switch (domain) {
    case "local":
      return `.${organizationSlug}.dev.starliner.cloud`;
    case "staging":
      return `.${organizationSlug}.staging.starliner.cloud`;
    case "production":
      return `.${organizationSlug}.starliner.cloud`;
  }
}

export function buildFullIngressHost(
  prefix: string,
  organizationSlug: string,
  domain: IngressHostDomain,
): string {
  const normalizedPrefix = sanitizeIngressHostPrefix(prefix);
  return `${normalizedPrefix}${getIngressHostSuffix(organizationSlug, domain)}`;
}

export function parseIngressHostPrefix(
  host: string,
  organizationSlug: string,
): string {
  const suffixes = [
    getIngressHostSuffix(organizationSlug, "production"),
    getIngressHostSuffix(organizationSlug, "staging"),
    getIngressHostSuffix(organizationSlug, "local"),
  ];

  for (const suffix of suffixes) {
    if (host.endsWith(suffix)) {
      return host.slice(0, -suffix.length);
    }
  }

  return host;
}

export function isValidIngressHostPrefix(prefix: string): boolean {
  const normalized = sanitizeIngressHostPrefix(prefix);
  if (!normalized) {
    return false;
  }

  return /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/.test(normalized);
}

export function sanitizeIngressHostPrefix(input: string): string {
  return input
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export function isValidIngressHost(
  host: string,
  organizationSlug: string,
  domain: IngressHostDomain,
): boolean {
  const suffix = getIngressHostSuffix(organizationSlug, domain);
  if (!host.endsWith(suffix)) {
    return false;
  }

  const prefix = host.slice(0, -suffix.length);
  return isValidIngressHostPrefix(prefix);
}

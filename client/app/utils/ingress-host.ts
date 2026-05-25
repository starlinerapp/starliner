const PREFIX_PATTERN = /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/;

function normalizePrefix(input: string): string {
  return input
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export function getIngressHostSuffix(
  organizationSlug: string,
  deploymentEnvironment: string,
): string {
  const subdomain =
    deploymentEnvironment === "local"
      ? "dev"
      : deploymentEnvironment === "staging"
        ? "staging"
        : null;

  return subdomain
    ? `.${organizationSlug}.${subdomain}.starliner.cloud`
    : `.${organizationSlug}.starliner.cloud`;
}

export function isValidIngressHostPrefix(prefix: string): boolean {
  const normalized = normalizePrefix(prefix);
  return normalized !== "" && PREFIX_PATTERN.test(normalized);
}

export function buildFullIngressHost(
  prefix: string,
  organizationSlug: string,
  deploymentEnvironment: string,
): string {
  return `${normalizePrefix(prefix)}${getIngressHostSuffix(organizationSlug, deploymentEnvironment)}`;
}

export function parseIngressHostPrefix(
  host: string,
  organizationSlug: string,
  deploymentEnvironment: string,
): string {
  const suffix = getIngressHostSuffix(organizationSlug, deploymentEnvironment);
  return host.endsWith(suffix) ? host.slice(0, -suffix.length) : host;
}

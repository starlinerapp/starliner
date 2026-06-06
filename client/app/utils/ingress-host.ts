const PREFIX_PATTERN = /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/;

export function getIngressHostSuffix(
  organizationSlug: string,
  deploymentEnvironment: string,
  deploymentDomain: string,
): string {
  const subdomain =
    deploymentEnvironment === "local"
      ? "dev"
      : deploymentEnvironment === "staging"
        ? "staging"
        : null;

  return subdomain
    ? `.${organizationSlug}.${subdomain}.${deploymentDomain}`
    : `.${organizationSlug}.${deploymentDomain}`;
}

export function isValidIngressHostPrefix(prefix: string): boolean {
  return PREFIX_PATTERN.test(prefix.trim().toLowerCase());
}

export function parseIngressHostPrefix(
  host: string,
  organizationSlug: string,
  deploymentEnvironment: string,
  deploymentDomain: string,
): string {
  const suffix = getIngressHostSuffix(
    organizationSlug,
    deploymentEnvironment,
    deploymentDomain,
  );
  return host.endsWith(suffix) ? host.slice(0, -suffix.length) : host;
}

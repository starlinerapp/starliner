package value

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidIngressHostPrefix = errors.New("invalid ingress host prefix")

var (
	ingressHostPrefixPattern      = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	ingressHostInvalidCharPattern = regexp.MustCompile(`[^a-z0-9-]+`)
	ingressHostDashPattern        = regexp.MustCompile(`-+`)
)

type IngressHostInput struct {
	Prefix string
	Paths  []*IngressPath
}

func NormalizeIngressHostPrefix(prefix string) string {
	normalized := strings.ToLower(strings.TrimSpace(prefix))
	normalized = ingressHostInvalidCharPattern.ReplaceAllString(normalized, "-")
	normalized = ingressHostDashPattern.ReplaceAllString(normalized, "-")
	normalized = strings.Trim(normalized, "-")
	return normalized
}

func IsValidIngressHostPrefix(prefix string) bool {
	normalized := NormalizeIngressHostPrefix(prefix)
	return normalized != "" && ingressHostPrefixPattern.MatchString(normalized)
}

func GetIngressHostSuffix(
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) string {
	subdomain := ""
	switch serverEnvironment {
	case "local":
		subdomain = "dev"
	case "staging":
		subdomain = "staging"
	}

	if subdomain != "" {
		return "." + organizationSlug + "." + subdomain + "." + deploymentDomain
	}

	return "." + organizationSlug + "." + deploymentDomain
}

func BuildFullIngressHost(
	prefix string,
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) string {
	return NormalizeIngressHostPrefix(prefix) + GetIngressHostSuffix(
		organizationSlug,
		serverEnvironment,
		deploymentDomain,
	)
}

func ParseIngressHostPrefix(
	host string,
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) string {
	suffix := GetIngressHostSuffix(organizationSlug, serverEnvironment, deploymentDomain)
	if strings.HasSuffix(host, suffix) {
		return host[:len(host)-len(suffix)]
	}
	return host
}

func BuildIngressHostsFromPrefixes(
	inputs []*IngressHostInput,
	organizationSlug string,
	serverEnvironment string,
	deploymentDomain string,
) ([]*IngressHost, error) {
	out := make([]*IngressHost, 0, len(inputs))

	for _, input := range inputs {
		if input == nil {
			continue
		}

		if !IsValidIngressHostPrefix(input.Prefix) {
			return nil, ErrInvalidIngressHostPrefix
		}

		out = append(out, &IngressHost{
			Host:  BuildFullIngressHost(input.Prefix, organizationSlug, serverEnvironment, deploymentDomain),
			Paths: input.Paths,
		})
	}

	return out, nil
}

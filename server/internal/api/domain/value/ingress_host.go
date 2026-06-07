package value

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidIngressHostPrefix = errors.New("invalid ingress host prefix")

var ingressHostPrefixPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)

type IngressHostInput struct {
	Prefix string
	Paths  []*IngressPath
}

func NormalizeIngressHostPrefix(prefix string) string {
	return strings.ToLower(strings.TrimSpace(prefix))
}

func IsValidIngressHostPrefix(prefix string) bool {
	return ingressHostPrefixPattern.MatchString(NormalizeIngressHostPrefix(prefix))
}

package value

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidIngressHostPrefix = errors.New("invalid ingress host prefix")

var ingressHostPrefixPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)

type IngressHostPrefix string

func NewIngressHostPrefix(prefix string) (IngressHostPrefix, error) {
	normalized := strings.ToLower(strings.TrimSpace(prefix))

	if !ingressHostPrefixPattern.MatchString(normalized) {
		return "", ErrInvalidIngressHostPrefix
	}

	return IngressHostPrefix(normalized), nil
}

type IngressHostInput struct {
	Prefix string
	Paths  []*IngressPath
}

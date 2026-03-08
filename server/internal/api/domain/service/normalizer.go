package service

import (
	"errors"
	"regexp"
	"strings"
)

type NormalizerService struct {
}

func NewNormalizerService() *NormalizerService {
	return &NormalizerService{}
}

func (ns *NormalizerService) FormatToDNS1123(label string) (string, error) {
	label = strings.ToLower(label)

	// Replace any character that isn't alphanumeric or hyphen with a hyphen
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	label = re.ReplaceAllString(label, "-")

	// Collapse multiple consecutive hyphens
	reHyphens := regexp.MustCompile(`-+`)
	label = reHyphens.ReplaceAllString(label, "-")

	// Trim leading and trailing hyphens
	label = strings.Trim(label, "-")

	if len(label) > 63 {
		label = label[:63]
		label = strings.TrimRight(label, "-")
	}

	if label == "" {
		return "", errors.New("invalid DNS label")
	}
	return label, nil
}

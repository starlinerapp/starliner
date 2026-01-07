package helm

import (
	"embed"
)

//go:embed template
var Chart embed.FS

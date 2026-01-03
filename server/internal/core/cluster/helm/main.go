package helm

import (
	"embed"
	_ "embed"
)

//go:embed template
var Chart embed.FS

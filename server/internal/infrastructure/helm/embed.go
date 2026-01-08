package helm

import (
	"embed"
)

//go:embed nginx
var NginxChart embed.FS

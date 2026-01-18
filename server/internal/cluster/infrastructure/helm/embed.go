package helm

import (
	"embed"
)

//go:embed postgres
var PostgresChart embed.FS

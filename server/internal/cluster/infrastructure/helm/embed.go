package helm

import (
	"embed"
)

//go:embed application
var ApplicationChart embed.FS

//go:embed cloudnative-pg
var CloudNativePgChart embed.FS

//go:embed postgres
var PostgresChart embed.FS

//go:embed ingress
var IngressChart embed.FS

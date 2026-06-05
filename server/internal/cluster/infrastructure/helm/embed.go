package helm

import (
	"embed"
)

//go:embed deployment
var DeploymentChart embed.FS

//go:embed statefulset
var StatefulSetChart embed.FS

//go:embed postgres
var PostgresChart embed.FS

//go:embed ingress
var IngressChart embed.FS

//go:embed external-dns
var ExternalDNSChart embed.FS

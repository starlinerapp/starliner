package mapper

import (
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
)

func MapHostsFromRequest(hosts []request.IngressHost) []*value.IngressHost {
	out := make([]*value.IngressHost, 0, len(hosts))
	for _, h := range hosts {
		paths := make([]*value.IngressPath, 0, len(h.Paths))
		for _, p := range h.Paths {
			paths = append(paths, &value.IngressPath{
				Path:        p.Path,
				PathType:    value.PathType(p.PathType),
				ServiceName: p.ServiceName,
			})
		}
		out = append(out, &value.IngressHost{
			Host:  h.Host,
			Paths: paths,
		})
	}
	return out
}

func MapEnvVarsFromRequest(envVars []request.EnvVar) []*value.EnvVar {
	out := make([]*value.EnvVar, 0, len(envVars))
	for _, e := range envVars {
		out = append(out, &value.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}
	return out
}

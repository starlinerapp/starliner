package service

import (
	"context"
	"fmt"
	"starliner.app/internal/api/domain/entity"
	interfaces "starliner.app/internal/api/domain/repository/interface"
	"strings"
)

type ResolverService struct {
	environmentRepository interfaces.EnvironmentRepository
}

func NewResolverService(
	environmentRepository interfaces.EnvironmentRepository,
) *ResolverService {
	return &ResolverService{
		environmentRepository: environmentRepository,
	}
}

func (rs *ResolverService) Resolve(ctx context.Context, environmentId int64, result ParseResult) (string, error) {
	if len(result.Literals) != len(result.Spans)+1 {
		return "", fmt.Errorf("malformed ParseResult")
	}

	var sb strings.Builder
	for i, span := range result.Spans {
		sb.WriteString(result.Literals[i])
		value, err := rs.resolveRef(ctx, environmentId, span.Ref)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintf(&sb, "%v", value)
		if err != nil {
			return "", err
		}
	}
	sb.WriteString(result.Literals[len(result.Literals)-1])
	return sb.String(), nil
}

func (rs *ResolverService) resolveRef(ctx context.Context, environmentId int64, ref TemplateRef) (any, error) {
	ingress, err := rs.environmentRepository.GetEnvironmentIngressDeploymentByName(ctx, environmentId, ref.Service)
	if err != nil {
		return ref.Raw, fmt.Errorf("failed to resolve service: %w", err)
	}

	var current any = ingress

	for _, step := range ref.Path {
		switch step.Kind {
		case StepKey:
			switch v := current.(type) {
			case *entity.IngressDeployment:
				switch step.Key {
				case "hosts":
					current = v.IngressHosts
				default:
					return ref.Raw, fmt.Errorf("unknown key %s for ingress", step.Key)
				}
			case *entity.IngressHost:
				switch step.Key {
				case "host":
					current = v.Host
				case "paths":
					current = v.Paths
				default:
					return ref.Raw, fmt.Errorf("unknown key %s for ingress host", step.Key)
				}
			}
		case StepIndex:
			switch v := current.(type) {
			case []*entity.IngressHost:
				if step.Index < 0 || step.Index >= len(v) {
					return ref.Raw, fmt.Errorf("hosts index %d out of range", step.Index)
				}
				current = v[step.Index]

			case []*entity.IngressPath:
				if step.Index < 0 || step.Index >= len(v) {
					return ref.Raw, fmt.Errorf("paths index %d out of range", step.Index)
				}
				current = v[step.Index]

			default:
				return ref.Raw, fmt.Errorf("cannot index type %T", current)
			}
		default:
			return ref.Raw, fmt.Errorf("unknown step kind %v", step.Kind)
		}
	}
	return current, nil
}

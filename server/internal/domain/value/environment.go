package value

import (
	"starliner.app/internal/domain/entity"
)

type Environment struct {
	Id   int64
	Slug string
	Name string
}

func NewEnvironment(e *entity.Environment) *Environment {
	return &Environment{
		Id:   e.Id,
		Slug: e.Slug,
		Name: e.Name,
	}
}

func NewEnvironments(es []*entity.Environment) []*Environment {
	environments := make([]*Environment, len(es))
	for i, e := range es {
		environments[i] = NewEnvironment(e)
	}
	return environments
}

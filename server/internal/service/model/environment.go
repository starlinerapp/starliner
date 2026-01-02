package model

import "starliner.app/internal/domain"

type Environment struct {
	Id   int64
	Slug string
	Name string
}

func NewEnvironment(e *domain.Environment) *Environment {
	return &Environment{
		Id:   e.Id,
		Slug: e.Slug,
		Name: e.Name,
	}
}

func NewEnvironments(es []*domain.Environment) []*Environment {
	environments := make([]*Environment, len(es))
	for i, e := range es {
		environments[i] = NewEnvironment(e)
	}
	return environments
}

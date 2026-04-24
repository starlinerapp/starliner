package response

import "starliner.app/internal/api/domain/value"

type Environment struct {
	Id   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

func NewEnvironment(environment *value.Environment) Environment {
	return Environment{
		Id:   environment.Id,
		Name: environment.Name,
		Slug: environment.Slug,
	}
}

func NewEnvironments(environments []*value.Environment) []Environment {
	result := make([]Environment, len(environments))
	for i, e := range environments {
		result[i] = Environment{
			Id:   e.Id,
			Name: e.Name,
			Slug: e.Slug,
		}
	}
	return result
}

type EnvironmentBranch struct {
	Branch string `json:"branch" binding:"required"`
}

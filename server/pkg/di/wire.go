//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "starliner.app/pkg/api"
	"starliner.app/pkg/api/handler"
	"starliner.app/pkg/config"
)

func InitializeAPI(cfg config.Config) (*http.Server, error) {
	wire.Build(
		http.NewServer,
		handler.NewRootHandler,
	)

	return &http.Server{}, nil
}

package api

import (
	"go.uber.org/fx"
	"starliner.app/pkg/api/handler"
	"starliner.app/pkg/api/middleware"
	"starliner.app/pkg/db"
	"starliner.app/pkg/repository"
	"starliner.app/pkg/service"
)

var Module = fx.Module(
	"api",
	db.Module,
	repository.Module,
	service.Module,
	middleware.Module,
	handler.Module,
	fx.Provide(NewServer),
	fx.Invoke(RegisterServer),
)

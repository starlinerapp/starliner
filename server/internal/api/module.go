package api

import (
	"go.uber.org/fx"
	"starliner.app/internal/api/handler"
	"starliner.app/internal/api/middleware"
	"starliner.app/internal/infrastructure/db"
	"starliner.app/internal/infrastructure/objectstore"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/repository"
	"starliner.app/internal/service"
)

var Module = fx.Module(
	"api",
	db.Module,
	objectstore.Module,
	queue.Module,
	repository.Module,
	service.Module,
	middleware.Module,
	handler.Module,
	fx.Provide(NewServer),
	fx.Invoke(RegisterServer),
)

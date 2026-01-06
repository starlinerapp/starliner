package http

import (
	"go.uber.org/fx"
	"starliner.app/internal/application"
	"starliner.app/internal/domain/repository"
	"starliner.app/internal/domain/service"
	"starliner.app/internal/infrastructure/crypto"
	"starliner.app/internal/infrastructure/dagger"
	"starliner.app/internal/infrastructure/db"
	"starliner.app/internal/infrastructure/objectstore"
	"starliner.app/internal/infrastructure/queue"
	"starliner.app/internal/presentation/http/handler"
	"starliner.app/internal/presentation/http/middleware"
)

var Module = fx.Module(
	"api",
	db.Module,
	objectstore.Module,
	queue.Module,
	repository.Module,
	dagger.Module,
	crypto.Module,
	service.Module,
	application.Module,
	middleware.Module,
	handler.Module,
	fx.Provide(NewServer),
	fx.Invoke(RegisterServer),
)

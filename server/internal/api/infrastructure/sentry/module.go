package sentry

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
	"starliner.app/internal/api/conf"
)

func InitSentry(lc fx.Lifecycle, cfg *conf.Config) error {
	if cfg.SentryDSN == "" {
		log.Println("sentry: DSN not set, skipping init")
		return nil
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.SentryDSN,
		Environment: cfg.Environment,
	}); err != nil {
		log.Printf("sentry: init failed, continuing without it: %v", err)
		return nil
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			sentry.Flush(2 * time.Second)
			return nil
		},
	})
	return nil
}

var Module = fx.Module(
	"sentry",
	fx.Invoke(InitSentry),
)

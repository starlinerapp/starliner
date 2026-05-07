package sentry

import (
	"context"
	"fmt"
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
		return fmt.Errorf("sentry init failed: %w", err)
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

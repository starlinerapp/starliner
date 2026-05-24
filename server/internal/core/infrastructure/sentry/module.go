package sentry

import (
	"context"
	"fmt"
	"log"
	"starliner.app/internal/core/conf"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

func InitSentry(lc fx.Lifecycle, sentryCfg conf.SentryConfig, envConfig conf.EnvironmentConfig) error {
	if sentryCfg.GetSentryDSN() == "" {
		log.Println("sentry: DSN not set, skipping init")
		return nil
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryCfg.GetSentryDSN(),
		Environment: envConfig.GetEnvironment(),
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

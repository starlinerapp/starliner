package concurrent

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func WithRecovery(ctx context.Context, name string, fn func() error) {
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s shutting down\n", name)
			return
		default:
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("worker", name)
						sentry.CurrentHub().Recover(r)
					})
					log.Printf("%s recovered from panic: %v\n", name, r)
				}
			}()

			if err := fn(); err != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("worker", name)
					sentry.CaptureException(err)
				})
				log.Printf("%s error: %v\n", name, err)
			}
		}()

		select {
		case <-ctx.Done():
			log.Printf("%s shutting down\n", name)
			return
		case <-time.After(backoff):
			log.Printf("%s restarting (backoff: %v)", name, backoff)
		}

		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

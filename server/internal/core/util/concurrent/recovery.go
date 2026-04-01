package concurrent

import (
	"context"
	"log"
	"time"
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
					log.Printf("%s recovered from panic: %v\n", name, r)
				}
			}()

			if err := fn(); err != nil {
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

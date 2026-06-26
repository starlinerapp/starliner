package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"starliner.app/internal/builder/domain/port"
	corePort "starliner.app/internal/core/domain/port"
)

const buildLogPollInterval = 100 * time.Millisecond

type BuildLogApplication struct {
	streams corePort.KVStore
}

var _ port.LogPublisher = (*BuildLogApplication)(nil)

func NewBuildLogApplication(streams corePort.KVStore) *BuildLogApplication {
	return &BuildLogApplication{
		streams: streams,
	}
}

func (a *BuildLogApplication) StreamBuildLogs(ctx context.Context, buildId int64) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	go func() {
		defer func() {
			_ = pw.Close()
		}()

		streamName := fmt.Sprintf("build:%d:logs", buildId)
		lastID := "0"

		for {
			select {
			case <-ctx.Done():
				_ = pw.CloseWithError(ctx.Err())
				return
			default:
			}

			entries, err := a.streams.ReadStream(ctx, streamName, lastID)
			if err != nil {
				if isBuildLogStreamNotReady(err) {
					if wait(ctx, buildLogPollInterval) {
						return
					}
					continue
				}
				_ = pw.CloseWithError(err)
				return
			}

			if len(entries) == 0 {
				if wait(ctx, buildLogPollInterval) {
					return
				}
				continue
			}

			for _, entry := range entries {
				lastID = entry.ID

				if _, ok := entry.Values["end"]; ok {
					return
				}

				data, ok := entry.Values["data"]
				if !ok || len(data) == 0 {
					continue
				}

				if _, err := pw.Write(data); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			}
		}
	}()

	return pr, nil
}

func (a *BuildLogApplication) PublishLogChunk(buildId int64, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	return a.streams.AppendToStream(
		context.Background(),
		fmt.Sprintf("build:%d:logs", buildId),
		map[string][]byte{
			"data": data,
		},
	)
}

func (a *BuildLogApplication) PublishLogEnd(buildId int64) error {
	ctx := context.Background()
	streamName := fmt.Sprintf("build:%d:logs", buildId)

	return a.streams.AppendToStream(ctx, streamName, map[string][]byte{
		"end": []byte("1"),
	})
}

func isBuildLogStreamNotReady(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, redis.Nil) {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "no such key") ||
		strings.Contains(msg, "unknown stream")
}

func wait(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return true
	case <-timer.C:
		return false
	}
}

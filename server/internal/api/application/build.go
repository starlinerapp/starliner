package application

import (
	"context"
	"io"

	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
)

type BuildApplication struct {
	buildRepository interfaces.BuildRepository
	pubsub          port.Pubsub
}

func NewBuildApplication(
	buildRepository interfaces.BuildRepository,
	pubsub port.Pubsub,
) *BuildApplication {
	return &BuildApplication{
		buildRepository: buildRepository,
		pubsub:          pubsub,
	}
}

func (ba *BuildApplication) GetBuildLogs(ctx context.Context, userId int64, buildId int64) (*string, error) {
	return ba.buildRepository.GetBuildLogs(ctx, userId, buildId)
}

func (ba *BuildApplication) StreamBuildLogs(
	ctx context.Context,
	userId int64,
	buildId int64,
	w io.Writer,
) error {
	// Subscribe first so that any chunks emitted while we verify access /
	// query the DB are not missed.
	chunks, cancelSub, err := ba.pubsub.SubscribeToBuildLogs(ctx, buildId)
	if err != nil {
		return err
	}
	defer cancelSub()

	logs, err := ba.buildRepository.GetBuildLogs(ctx, userId, buildId)
	if err != nil {
		return err
	}

	// If a completed build already has logs persisted, send those.
	if logs != nil && *logs != "" {
		if _, err := io.WriteString(w, *logs); err != nil {
			return err
		}
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunk, ok := <-chunks:
			if !ok {
				return nil
			}
			if chunk.End {
				return nil
			}
			if len(chunk.Data) == 0 {
				continue
			}
			if _, err := w.Write(chunk.Data); err != nil {
				return err
			}
		}
	}
}

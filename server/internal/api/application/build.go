package application

import (
	"context"
	"errors"
	"io"

	"starliner.app/internal/api/domain/port"
	interfaces "starliner.app/internal/api/domain/repository/interface"
)

type BuildApplication struct {
	buildRepository   interfaces.BuildRepository
	grpcBuilderClient port.BuilderClient
}

func NewBuildApplication(
	buildRepository interfaces.BuildRepository,
	grpcBuilderClient port.BuilderClient,
) *BuildApplication {
	return &BuildApplication{
		buildRepository:   buildRepository,
		grpcBuilderClient: grpcBuilderClient,
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
	pr, pw := io.Pipe()
	streamCtx, cancelStream := context.WithCancel(ctx)
	defer cancelStream()

	errCh := make(chan error, 1)
	go func() {
		errCh <- ba.grpcBuilderClient.StreamBuildLogs(streamCtx, buildId, pw)
		_ = pw.Close()
	}()

	logs, err := ba.buildRepository.GetBuildLogs(ctx, userId, buildId)
	if err != nil {
		cancelStream()
		<-errCh
		_ = pr.Close()
		return err
	}

	// A completed build may already have full logs in the database.
	if logs != nil && *logs != "" {
		cancelStream()
		<-errCh
		_ = pr.Close()
		_, werr := io.WriteString(w, *logs)
		return werr
	}

	_, copyErr := io.Copy(w, pr)
	grpcErr := <-errCh
	if copyErr != nil {
		return copyErr
	}
	if grpcErr != nil && !errors.Is(grpcErr, context.Canceled) {
		return grpcErr
	}
	return nil
}

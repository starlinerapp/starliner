package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToBuildTriggered(handler func(build *value.TriggerBuild)) error
	PublishBuildCompleted(build *value.BuildCompleted) error
	PublishRunnerStatusChanged(status *value.RunnerStatusChanged) error
	SubscribeToRunnerDeleted(handler func(runner *value.RunnerDeleted)) error
}

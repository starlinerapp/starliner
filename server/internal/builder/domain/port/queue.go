package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToBuildTriggered(handler func(build *value.TriggerBuild)) error
	PublishBuildCompleted(build *value.BuildCompleted) error
	PublishRunnerStatusChanged(status *value.RunnerStatusChanged) error
	SubscribeToRunnerDeleted(handler func(runner *value.RunnerDeleted)) error

	PublishRunnerJob(job *value.RunnerBuildJob) error
	ClaimRunnerJob(runnerId int64) (*value.RunnerBuildJob, error)
	PublishRunnerJobResult(result *value.RunnerBuildResult) error
	SubscribeToRunnerJobResults(handler func(result *value.RunnerBuildResult)) error
}

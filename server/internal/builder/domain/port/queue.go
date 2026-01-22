package port

import "starliner.app/internal/core/domain/value"

type Queue interface {
	SubscribeToBuildTriggered(handler func(build *value.TriggerBuild)) error
}

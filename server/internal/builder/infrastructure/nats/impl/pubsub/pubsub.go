package pubsub

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nats-io/nats.go"
	"starliner.app/internal/builder/domain/port"
	"starliner.app/internal/core/domain/value"
	natscore "starliner.app/internal/core/infrastructure/nats/core"
)

// BuildLogs is the NATS subject prefix used to stream build log chunks and the
// terminal end marker. The full subject is "build.logs.<buildId>".
const BuildLogs natscore.Subject = "build.logs"

type Pubsub struct {
	publisher *natscore.Publisher
}

func NewPubsub(conn *nats.Conn) port.LogPublisher {
	return &Pubsub{
		publisher: natscore.NewPublisher(conn),
	}
}

func (p *Pubsub) PublishLogChunk(buildId int64, chunk []byte) error {
	data := make([]byte, len(chunk))
	copy(data, chunk)

	payload, err := json.Marshal(&value.BuildLogChunk{
		BuildId: buildId,
		Data:    data,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal log chunk: %w", err)
	}
	return p.publisher.Publish(BuildLogs, strconv.FormatInt(buildId, 10), payload)
}

func (p *Pubsub) PublishLogEnd(buildId int64) error {
	payload, err := json.Marshal(&value.BuildLogChunk{
		BuildId: buildId,
		End:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal log end: %w", err)
	}
	return p.publisher.Publish(BuildLogs, strconv.FormatInt(buildId, 10), payload)
}

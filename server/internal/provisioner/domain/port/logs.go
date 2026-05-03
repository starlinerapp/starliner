package port

type LogPublisher interface {
	PublishLogChunk(clusterId int64, chunk []byte) error
	PublishLogEnd(clusterId int64) error
}

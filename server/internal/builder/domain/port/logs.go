package port

type LogPublisher interface {
	PublishLogChunk(buildId int64, chunk []byte) error
	PublishLogEnd(buildId int64) error
}

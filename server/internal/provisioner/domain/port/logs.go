package port

type LogPublisher interface {
	PublishLogChunk(provisioningId string, chunk []byte) error
	PublishLogEnd(provisioningId string) error
}

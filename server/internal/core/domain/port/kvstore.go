package port

import "context"

type StreamEntry struct {
	ID     string
	Values map[string][]byte
}

type KVStore interface {
	AppendToStream(ctx context.Context, name string, payload map[string][]byte) error
	ReadStream(ctx context.Context, name string, lastId string) ([]StreamEntry, error)
}

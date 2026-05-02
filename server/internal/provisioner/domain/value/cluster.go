package value

type ClusterLogChunk struct {
	ClusterId int64
	Data      []byte
	End       bool
}

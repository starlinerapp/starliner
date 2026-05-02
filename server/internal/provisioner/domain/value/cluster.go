package value

type ClusterLogChunk struct {
	ProvisioningId string
	Data           []byte
	End            bool
}

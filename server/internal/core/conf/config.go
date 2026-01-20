package conf

type S3Config interface {
	GetS3EndpointUrl() string
	GetAWSAccessKeyId() string
	GetAWSSecretAccessKey() string
}

type NatsConfig interface {
	GetNatsUrl() string
}

type CryptoConfig interface {
	GetEncryptionKeyBase64() string
}

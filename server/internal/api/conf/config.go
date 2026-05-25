package conf

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost                  string `mapstructure:"DB_HOST" validate:"required"`
	DBPort                  int    `mapstructure:"DB_PORT" validate:"required"`
	DBUser                  string `mapstructure:"DB_USER" validate:"required"`
	DBPassword              string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBName                  string `mapstructure:"DB_NAME" validate:"required"`
	BasicAuthUser           string `mapstructure:"BASIC_AUTH_USER" validate:"required"`
	BasicAuthPassword       string `mapstructure:"BASIC_AUTH_PASSWORD" validate:"required"`
	ClusterGrpcEndpoint     string `mapstructure:"CLUSTER_GRPC_ENDPOINT" validate:"required"`
	BuilderGrpcEndpoint     string `mapstructure:"BUILDER_GRPC_ENDPOINT" validate:"required"`
	ProvisionerGrpcEndpoint string `mapstructure:"PROVISIONER_GRPC_ENDPOINT" validate:"required"`
	S3EndpointUrl           string `mapstructure:"S3_ENDPOINT_URL" validate:"required"`
	NatsUrl                 string `mapstructure:"NATS_URL" validate:"required"`
	AWSAccessKeyId          string `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey      string `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	EncryptionKeyBase64     string `mapstructure:"ENCRYPTION_KEY_BASE64" validate:"required"`
	GithubAppPrivateKey     string `mapstructure:"GITHUB_APP_PRIVATE_KEY" validate:"required"`
	GithubAppID             int64  `mapstructure:"GITHUB_APP_ID" validate:"required"`
	GithubWebhookSecret     string `mapstructure:"GITHUB_WEBHOOK_SECRET" validate:"required"`
	SenderMail              string `mapstructure:"SENDER_MAIL" validate:"required"`
	SmtpHost                string `mapstructure:"SMTP_HOST" validate:"required"`
	SmtpPort                string `mapstructure:"SMTP_PORT" validate:"required"`
	SmtpUsername            string `mapstructure:"SMTP_USERNAME"`
	SmtpPassword            string `mapstructure:"SMTP_PASSWORD"`
	SmtpTLSEnabled          bool   `mapstructure:"SMTP_TLS_ENABLED"`
	ImageRegistryUrl        string `mapstructure:"IMAGE_REGISTRY_URL" validate:"required"`
	ImageRegistryUsername   string `mapstructure:"IMAGE_REGISTRY_USERNAME" validate:"required"`
	ImageRegistryPassword   string `mapstructure:"IMAGE_REGISTRY_PASSWORD" validate:"required"`
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfgType := reflect.TypeOf(config)
	for i := 0; i < cfgType.NumField(); i++ {
		if env := cfgType.Field(i).Tag.Get("mapstructure"); env != "" {
			if err := viper.BindEnv(env); err != nil {
				return nil, err
			}
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	if err := validator.New().Struct(config); err != nil {
		return &config, err
	}

	return &config, nil
}

func (c *Config) GetS3EndpointUrl() string {
	return c.S3EndpointUrl
}

func (c *Config) GetAWSAccessKeyId() string {
	return c.AWSAccessKeyId
}

func (c *Config) GetAWSSecretAccessKey() string {
	return c.AWSSecretAccessKey
}

func (c *Config) GetNatsUrl() string {
	return c.NatsUrl
}

func (c *Config) GetEncryptionKeyBase64() string {
	return c.EncryptionKeyBase64
}

package conf

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"reflect"
)

type Config struct {
	S3EndpointUrl      string `mapstructure:"S3_ENDPOINT_URL" validate:"required"`
	NatsUrl            string `mapstructure:"NATS_URL" validate:"required"`
	AWSAccessKeyId     string `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	ImageRegistryUrl   string `mapstructure:"IMAGE_REGISTRY_URL" validate:"required"`
	BuilderGrpcAddr    string `mapstructure:"BUILDER_GRPC_ADDR" validate:"required"`
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

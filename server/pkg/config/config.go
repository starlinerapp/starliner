package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost             string `mapstructure:"DB_HOST" validate:"required"`
	DBPort             int    `mapstructure:"DB_PORT" validate:"required"`
	DBUser             string `mapstructure:"DB_USER" validate:"required"`
	DBPassword         string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBName             string `mapstructure:"DB_NAME" validate:"required"`
	BasicAuthUser      string `mapstructure:"BASIC_AUTH_USER" validate:"required"`
	BasicAuthPassword  string `mapstructure:"BASIC_AUTH_PASSWORD" validate:"required"`
	S3EndpointUrl      string `mapstructure:"S3_ENDPOINT_URL" validate:"required"`
	NatsUrl            string `mapstructure:"NATS_URL" validate:"required"`
	AWSAccessKeyId     string `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	ImageRegistryUrl   string `mapstructure:"IMAGE_REGISTRY_URL" validate:"required"`
}

var envs = []string{
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
	"BASIC_AUTH_USER",
	"BASIC_AUTH_PASSWORD",
	"S3_ENDPOINT_URL",
	"NATS_URL",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"IMAGE_REGISTRY_URL",
}

func LoadConfig() (*Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return &config, err
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

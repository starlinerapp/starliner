package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"starliner.app/internal/core/conf"
)

func Connect(cfg conf.S3Config) (*s3.Client, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.GetAWSAccessKeyId(),
				SecretAccessKey: cfg.GetAWSSecretAccessKey(),
			},
		}),
	)

	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.GetS3EndpointUrl())
		o.Region = "none"
		o.UsePathStyle = true
	})

	return client, nil
}

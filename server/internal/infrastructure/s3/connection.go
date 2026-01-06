package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"starliner.app/internal/conf"
)

func Connect(cfg *conf.Config) (*s3.Client, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.AWSAccessKeyId,
				SecretAccessKey: cfg.AWSSecretAccessKey,
			},
		}),
	)

	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3EndpointUrl)
		o.Region = "none"
		o.UsePathStyle = true
	})

	return client, nil
}

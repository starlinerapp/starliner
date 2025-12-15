package objectstore

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"starliner.app/pkg/config"
)

func Connect(cfg *config.Config) (*s3.Client, error) {
	log.Println("Connecting to S3")
	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.AWSAccessKeyId,
				SecretAccessKey: cfg.AWSSecretAccessKey,
			},
		}),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://s3:8333")
		o.Region = "us-east-1"
		o.UsePathStyle = true
	})

	return client, nil
}

func CreateBucket(client *s3.Client) {
	log.Println("Create bucket")
	createBucketParams := &s3.CreateBucketInput{
		Bucket: aws.String("data"),
	}
	_, err := client.CreateBucket(context.TODO(), createBucketParams)
	if err != nil {
		log.Println(fmt.Errorf("failed to create bucket: %v", err))
	}
}

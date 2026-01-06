package s3

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"log"
	"starliner.app/internal/domain/port"
)

const AppBucketName = "data"
const PulumiBucketName = "pulumi-backend"

type S3Client struct {
	client *s3.Client
}

func NewS3Client(client *s3.Client) port.ObjectStore {
	return &S3Client{client: client}
}

func (c *S3Client) CreateBuckets(ctx context.Context) error {
	buckets := []string{AppBucketName, PulumiBucketName}

	for _, bucketName := range buckets {
		input := &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		}

		_, err := c.client.CreateBucket(ctx, input)
		if err != nil {
			var bExists *types.BucketAlreadyExists
			if errors.As(err, &bExists) {
				continue
			}
			log.Printf("failed to create bucket %s: %v", bucketName, err)
			return err
		}

		log.Printf("Bucket %s created successfully", bucketName)
	}
	return nil
}
func (c *S3Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	res, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(AppBucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

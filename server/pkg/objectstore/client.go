package objectstore

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"log"
)

const BucketName = "data"

type S3Client struct {
	client *s3.Client
}

func NewS3Client(client *s3.Client) *S3Client {
	return &S3Client{client: client}
}

func (c *S3Client) CreateBucket(ctx context.Context) error {
	createBucketParams := &s3.CreateBucketInput{
		Bucket: aws.String(BucketName),
	}
	_, err := c.client.CreateBucket(ctx, createBucketParams)
	if err != nil {
		var bExists *types.BucketAlreadyExists
		if errors.As(err, &bExists) {
			return nil
		}

		log.Printf("failed to create bucket: %v", err)
		return err
	}
	return nil
}

func (c *S3Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	res, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

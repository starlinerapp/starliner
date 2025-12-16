package objectstore

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/fx"
	"log"
)

const BucketName = "data"

func CreateBucket(client *s3.Client, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: aws.String(BucketName),
			})

			if err != nil {
				var bExists *types.BucketAlreadyExists
				if errors.As(err, &bExists) {
					return nil
				}

				log.Printf("failed to create bucket: %v", err)
				return err
			}
			return nil
		},
	})
}

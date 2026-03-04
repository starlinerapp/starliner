package s3

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"starliner.app/internal/core/domain/port"
	"strings"
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

func (c *S3Client) UploadDirAsTarGz(ctx context.Context, pathToDir string, key string) (string, error) {
	preader, pwriter := io.Pipe()
	go func() {
		writeErr := c.writeTarGz(pwriter, pathToDir)
		_ = pwriter.CloseWithError(writeErr)
	}()

	uploader := transfermanager.New(c.client)

	_, err := uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket:          aws.String(AppBucketName),
		Key:             aws.String(key),
		Body:            preader,
		ContentType:     aws.String("application/gzip"),
		ContentEncoding: aws.String("gzip"),
	})

	if err != nil {
		return "", fmt.Errorf("upload to s3 failed: %w", err)
	}

	return key, nil
}

func (c *S3Client) writeTarGz(w io.Writer, pathToDir string) error {
	gzwriter := gzip.NewWriter(w)
	defer func(gzwriter *gzip.Writer) {
		err := gzwriter.Close()
		if err != nil {
			log.Printf("failed to close gzip writer: %v", err)
		}
	}(gzwriter)

	tw := tar.NewWriter(gzwriter)
	defer func(tw *tar.Writer) {
		err := tw.Close()
		if err != nil {
			log.Printf("failed to close tar writer: %v", err)
		}
	}(tw)

	root := filepath.Clean(pathToDir)
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		var link string
		if info.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}

		hdr.Name = filepath.ToSlash(rel)

		if info.IsDir() && !strings.HasSuffix(hdr.Name, "/") {
			hdr.Name += "/"
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(tw, file)
		closeErr := file.Close()
		if err != nil {
			return err
		}

		return closeErr
	})
}

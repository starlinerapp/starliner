package dagger

import (
	"context"
	"dagger.io/dagger"
	"os"
)

func NewDaggerClient() (*dagger.Client, error) {
	ctx := context.Background()
	return dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
}

package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"starliner.app/internal/core/domain/port"
)

const (
	monitoredRunnersKey = "runners:monitored"
	runnerStatusKeyFmt  = "runner:%d:status"
	runnerDeletedKeyFmt = "runner:%d:deleted"
)

type Client struct {
	client *redis.Client
}

func NewClient(client *redis.Client) *Client {
	return &Client{client: client}
}

func (c *Client) TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, "1", ttl).Result()
}

func (c *Client) AppendToStream(ctx context.Context, name string, payload map[string][]byte) error {
	values := make(map[string]any, len(payload))

	for key, value := range payload {
		values[key] = value
	}

	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: name,
		MaxLen: 10_000,
		Approx: true,
		Values: values,
	}).Err()
}

func (c *Client) ReadStream(ctx context.Context, name string, lastId string) ([]port.StreamEntry, error) {
	if lastId == "" {
		lastId = "0"
	}

	streams, err := c.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{name, lastId},
		Count:   100,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]port.StreamEntry, 0)

	for _, stream := range streams {
		for _, message := range stream.Messages {
			values := make(map[string][]byte, len(message.Values))

			for key, value := range message.Values {
				switch v := value.(type) {
				case string:
					values[key] = []byte(v)
				case []byte:
					values[key] = v
				}
			}

			entries = append(entries, port.StreamEntry{
				ID:     message.ID,
				Values: values,
			})
		}
	}

	return entries, nil
}

func (c *Client) MarkAlive(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Set(ctx, fmt.Sprintf("live:%s", key), "1", ttl).Err()
}

func (c *Client) IsAlive(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, fmt.Sprintf("live:%s", key)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Client) AddMonitoredRunner(ctx context.Context, runnerId int64) error {
	return c.client.SAdd(ctx, monitoredRunnersKey, runnerId).Err()
}

func (c *Client) RemoveMonitoredRunner(ctx context.Context, runnerId int64) error {
	pipe := c.client.Pipeline()
	pipe.SRem(ctx, monitoredRunnersKey, runnerId)
	pipe.Del(ctx, fmt.Sprintf(runnerStatusKeyFmt, runnerId))
	pipe.Del(ctx, fmt.Sprintf("live:runner:%d", runnerId))
	pipe.Set(ctx, fmt.Sprintf(runnerDeletedKeyFmt, runnerId), "1", 0)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *Client) IsRunnerDeleted(ctx context.Context, runnerId int64) (bool, error) {
	count, err := c.client.Exists(ctx, fmt.Sprintf(runnerDeletedKeyFmt, runnerId)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Client) ListMonitoredRunners(ctx context.Context) ([]int64, error) {
	members, err := c.client.SMembers(ctx, monitoredRunnersKey).Result()
	if err != nil {
		return nil, err
	}

	runnerIds := make([]int64, 0, len(members))
	for _, member := range members {
		runnerId, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			continue
		}
		runnerIds = append(runnerIds, runnerId)
	}

	return runnerIds, nil
}

func (c *Client) GetRunnerStatus(ctx context.Context, runnerId int64) (string, error) {
	status, err := c.client.Get(ctx, fmt.Sprintf(runnerStatusKeyFmt, runnerId)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return status, nil
}

func (c *Client) SetRunnerStatus(ctx context.Context, runnerId int64, status string) error {
	return c.client.Set(ctx, fmt.Sprintf(runnerStatusKeyFmt, runnerId), status, 0).Err()
}

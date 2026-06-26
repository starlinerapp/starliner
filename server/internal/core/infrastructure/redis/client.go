package redis

import "github.com/redis/go-redis/v9"

type Client struct {
	client *redis.Client
}

func NewClient(client *redis.Client) *Client {
	return &Client{client: client}
}

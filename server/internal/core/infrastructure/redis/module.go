package redis

import (
	"go.uber.org/fx"
	"starliner.app/internal/core/domain/port"
)

var Module = fx.Module(
	"redis",
	fx.Provide(
		Connect,
		NewClient,
		func(c *Client) port.KVStore { return c },
		func(c *Client) port.AcquireLimiter { return c },
	),
)

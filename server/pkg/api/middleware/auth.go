package middleware

import (
	"github.com/gin-gonic/gin"
	"starliner.app/pkg/config"
)

type BasicAuthMiddleware struct {
	cfg *config.Config
}

func NewBasicAuthMiddleware(cfg *config.Config) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{cfg: cfg}
}

func (b *BasicAuthMiddleware) WithBasicAuth() gin.HandlerFunc {
	user := b.cfg.BasicAuthUser
	password := b.cfg.BasicAuthPassword

	return gin.BasicAuth(gin.Accounts{
		user: password,
	})
}

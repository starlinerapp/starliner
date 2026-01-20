package middleware

import (
	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/conf"
)

type BasicAuthMiddleware struct {
	cfg *conf.Config
}

func NewBasicAuthMiddleware(cfg *conf.Config) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{cfg: cfg}
}

func (b *BasicAuthMiddleware) WithBasicAuth() gin.HandlerFunc {
	user := b.cfg.BasicAuthUser
	password := b.cfg.BasicAuthPassword

	return gin.BasicAuth(gin.Accounts{
		user: password,
	})
}

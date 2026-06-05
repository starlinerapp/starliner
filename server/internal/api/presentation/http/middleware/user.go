package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"starliner.app/internal/api/application"
)

type UserMiddleware struct {
	userApplication *application.UserApplication
}

func NewUserMiddleware(userApplication *application.UserApplication) *UserMiddleware {
	return &UserMiddleware{userApplication: userApplication}
}

func (u *UserMiddleware) WithUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("X-User-ID")
		if userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := u.userApplication.GetOrCreateUser(c, userId)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

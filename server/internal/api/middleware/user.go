package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"starliner.app/internal/service"
)

type UserMiddleware struct {
	userService *service.UserService
}

func NewUserMiddleware(userService *service.UserService) *UserMiddleware {
	return &UserMiddleware{userService: userService}
}

func (u *UserMiddleware) WithUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("X-User-ID")
		if userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := u.userService.GetOrCreateUser(c, userId)
		if err != nil {
			log.Printf("failed to get or create user")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

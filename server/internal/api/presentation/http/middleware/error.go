package middleware

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// WithErrorReporting runs the handler chain, then reports any errors attached
// via c.Error(...) to Sentry through the request hub. If the handler aborted
// without writing a response, it writes a standard 500.
func WithErrorReporting() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			for _, e := range c.Errors {
				hub.CaptureException(e.Err)
			}
		}

		if !c.Writer.Written() {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	}
}

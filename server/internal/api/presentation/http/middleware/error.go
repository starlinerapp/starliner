package middleware

import (
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// WithErrorReporting runs the handler chain, then reports any errors attached
// via c.Error(...) to Sentry through the request hub.
func WithErrorReporting() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		if c.Writer.Status() < http.StatusInternalServerError {
			return
		}

		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			for _, e := range c.Errors {
				hub.CaptureException(e.Err)
			}
		}
	}
}

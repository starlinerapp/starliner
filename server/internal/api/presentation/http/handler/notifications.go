package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/presentation/http/sse"
)

type NotificationsHandler struct {
	notificationApplication *application.NotificationApplication
}

func NewNotificationsHandler(
	notificationApplication *application.NotificationApplication,
) *NotificationsHandler {
	return &NotificationsHandler{
		notificationApplication: notificationApplication,
	}
}

// StreamGlobalNotifications godoc
// @Summary Stream global notifications
// @State core
// @Tags notifications
// @ID streamGlobalNotifications
// @Param X-User-ID header string true "User ID"
// @Param organizationId query int true "Organization ID"
// @Product text/event-stream
// @Success 200
// @Header 200 {string} Content-Type "text/event-stream"
// @Header 200 {string} Cache-Control "no-cache"
// @Header 200 {string} Connection "keep-alive"
// @Router /notifications [get]
func (nh *NotificationsHandler) StreamGlobalNotifications(c *gin.Context) {
	correlationId := c.GetHeader("X-Correlation-ID")
	if correlationId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Correlation-ID header"})
		return
	}

	_, err := strconv.ParseInt(c.Query("organizationId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid organizationId query parameter"})
		return
	}

	sw, ok := sse.NewWriter(c.Writer)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	sub := nh.notificationApplication.SubscribeGlobal(correlationId)
	defer sub.Close()

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case notification := <-sub.Notifications():
			sw.WriteJSON(notification)
		}
	}
}

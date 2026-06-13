package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/sse"
)

type NotificationsHandler struct {
	userNotificationHub *sse.UserNotificationHub
}

func NewNotificationsHandler(
	userNotificationHub *sse.UserNotificationHub,
) *NotificationsHandler {
	return &NotificationsHandler{
		userNotificationHub: userNotificationHub,
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
	currentUser := c.MustGet("user").(*value.User)

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

	ch := nh.userNotificationHub.Subscribe(currentUser.Id)
	defer nh.userNotificationHub.Unsubscribe(currentUser.Id, ch)

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case notification := <-ch:
			nh.userNotificationHub.WriteNotification(sw, notification)
		}
	}
}

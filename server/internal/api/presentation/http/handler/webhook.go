package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"starliner.app/internal/api/application"
)

type WebhookHandler struct {
	githubApplication *application.GitHubApplication
}

func NewWebhookHandler(
	githubApplication *application.GitHubApplication,
) *WebhookHandler {
	return &WebhookHandler{
		githubApplication: githubApplication,
	}
}

func (wh *WebhookHandler) HandleGithubWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	eventType := c.GetHeader("X-GitHub-Event")

	err = wh.githubApplication.HandleGithubWebhook(c.Request.Context(), eventType, payload)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

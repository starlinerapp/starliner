package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"starliner.app/internal/api/application"
	"starliner.app/internal/api/domain/port"
	"starliner.app/internal/api/domain/value"
	"starliner.app/internal/api/presentation/http/dto/request"
	"starliner.app/internal/api/presentation/http/dto/response"
	"starliner.app/internal/api/presentation/http/sse"
	"strconv"
)

type ClusterHandler struct {
	clusterApplication      *application.ClusterApplication
	organizationApplication *application.OrganizationApplication
}

func NewClusterHandler(clusterApplication *application.ClusterApplication, organizationApplication *application.OrganizationApplication) *ClusterHandler {
	return &ClusterHandler{
		clusterApplication:      clusterApplication,
		organizationApplication: organizationApplication,
	}
}

// CreateCluster FindAll godoc
// @Summary Create Cluster
// @Tags cluster
// @ID createCluster
// @Param X-User-ID header string true "User ID"
// @Param data body request.CreateCluster true "Create Cluster"
// @Product JSON
// @Success 200 {object} response.Cluster
// @Router /clusters [post]
func (ch *ClusterHandler) CreateCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	var cluster request.CreateCluster
	if err := c.BindJSON(&cluster); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	newCluster, err := ch.clusterApplication.CreateCluster(c.Request.Context(), currentUser.Id, cluster.Name, cluster.ServerType, cluster.OrganizationID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, response.NewCluster(newCluster))
}

// GetCluster FindAll godoc
// @Summary Get Cluster
// @Tags cluster
// @ID getCluster
// @Param X-User-ID header string true "User ID"
// @Product JSON
// @Param id path int true "Cluster ID"
// @Success 200 {object} response.Cluster
// @Router /clusters/{id} [get]
func (ch *ClusterHandler) GetCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cluster, err := ch.clusterApplication.GetUserCluster(c.Request.Context(), currentUser.Id, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	if cluster == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}
	c.JSON(http.StatusOK, response.NewCluster(cluster))
}

// GetClusterPrivateKey FindAll godoc
// @Summary Get Cluster Private Key
// @Tags cluster
// @ID getClusterPrivateKey
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Cluster ID"
// @Product application/octet-stream
// @Success 200 {file} file
// @Router /clusters/{id}/private-key [get]
func (ch *ClusterHandler) GetClusterPrivateKey(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	file, err := ch.clusterApplication.GetClusterPrivateKey(c.Request.Context(), clusterId, currentUser.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=private-key.pem")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(file)))

	c.Data(http.StatusOK, "application/octet-stream", file)
}

// DeleteCluster FindAll godoc
// @Summary Delete Cluster
// @Tags cluster
// @ID deleteCluster
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Cluster ID"
// @Product JSON
// @Success 200
// @Router /clusters/{id} [delete]
func (ch *ClusterHandler) DeleteCluster(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = ch.clusterApplication.DeleteCluster(c.Request.Context(), currentUser.Id, clusterId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	c.Status(http.StatusOK)
}

// StreamProvisioningLogs FindAll godoc
// @Summary Stream cluster provisioning logs
// @Tags cluster
// @ID streamClusterProvisioningLogs
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Cluster ID"
// @Product JSON
// @Success 200
// @Header 200 {string} Content-Type "text/event-stream"
// @Header 200 {string} Cache-Control "no-cache"
// @Header 200 {string} Connection "keep-alive"
// @Router /clusters/{id}/provisioning/logs/stream [get]
func (ch *ClusterHandler) StreamProvisioningLogs(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
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

	err = ch.clusterApplication.StreamProvisioningLogs(c.Request.Context(), currentUser.Id, clusterId, sw)
	if err != nil {
		sw.WriteError(err)
	}
}

var clusterUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ch *ClusterHandler) OpenTTY(c *gin.Context) {
	currentUser := c.MustGet("user").(*value.User)
	clusterId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	rows, _ := strconv.Atoi(c.Query("tty_height"))
	cols, _ := strconv.Atoi(c.Query("tty_width"))

	conn, err := clusterUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "websocket upgrade failed"})
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	sizeCh := make(chan port.TerminalSize, 1)
	defer close(sizeCh)

	sizeCh <- port.TerminalSize{Rows: rows, Columns: cols}

	go func() {
		defer func(stdinWriter *io.PipeWriter) {
			_ = stdinWriter.Close()
		}(stdinWriter)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if _, err := stdinWriter.Write(msg); err != nil {
				return
			}
		}
	}()

	go func() {
		defer func(stdoutReader *io.PipeReader) {
			_ = stdoutReader.Close()
		}(stdoutReader)
		buf := make([]byte, 4096)
		for {
			n, err := stdoutReader.Read(buf)
			if n > 0 {
				_ = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	err = ch.clusterApplication.OpenTTY(c.Request.Context(), currentUser.Id, clusterId, stdinReader, stdoutWriter, sizeCh)
	if err != nil && !errors.Is(err, io.EOF) {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("error: %v", err)))
	}
}

package handler

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"starliner.app/internal/api/domain/port"
)

const ttyResizeMessageType = "resize"

type ttyResizeMessage struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

func parseTTYResizeMessage(messageType int, msg []byte) (port.TerminalSize, bool) {
	if messageType != websocket.TextMessage {
		return port.TerminalSize{}, false
	}

	var resize ttyResizeMessage
	if err := json.Unmarshal(msg, &resize); err != nil || resize.Type != ttyResizeMessageType {
		return port.TerminalSize{}, false
	}
	if resize.Cols <= 0 || resize.Rows <= 0 {
		return port.TerminalSize{}, false
	}

	return port.TerminalSize{Rows: resize.Rows, Columns: resize.Cols}, true
}

func pushTerminalSize(sizeCh chan<- port.TerminalSize, size port.TerminalSize) {
	sizeCh <- size
}

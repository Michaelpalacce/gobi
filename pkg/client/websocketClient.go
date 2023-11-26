package client

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/gorilla/websocket"
)

// WebsocketClient contains the connection as well as metadata for a client
// This is mainly a transport layer connection
type WebsocketClient struct {
	Conn   *websocket.Conn
	Client Client
}

// Close will gracefully close the connection. If an error ocurrs during closing, it will be ignored.
func (c *WebsocketClient) Close(msg string) {
	payload := messages.NewCloseMessage(msg)

	// Close the WebSocket connection gracefully
	_ = c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, string(payload.Marshal())))
}

// sendMessage enforces a uniform style in sending data
func (c *WebsocketClient) SendMessage(message messages.WebsocketRequest) error {
	messageBytes := message.Marshal()
	slog.Debug("Sending message", "message", string(messageBytes))
	err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes)

	if err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

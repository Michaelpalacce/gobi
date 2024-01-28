package socket

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/Michaelpalacce/gobi/pkg/storage"
	"github.com/gorilla/websocket"
)

// WebsocketClient contains the connection as well as metadata for a client
// Used by both the server and client
// This is mainly a transport layer connection
type WebsocketClient struct {
	// General
	Conn          *websocket.Conn
	StorageDriver storage.Driver
	Client        client.Client
	closed        bool
	InitialSync   bool

	// Server Exclusive
	DB *database.Database

	// Client Exclusive
}

// Close will gracefully close the connection. If an error ocurrs during closing, it will be ignored.
// It will set the WebsocketClient as closed and will NOT send a CLose Message if the connection is closed already
func (c *WebsocketClient) Close(msg string) {
	if c.closed {
		return
	}

	c.closed = true
	payload := messages.NewCloseMessage(msg)

	// Close the WebSocket connection gracefully
	_ = c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, string(payload.Marshal())))
}

// sendMessage enforces a uniform style in sending data
func (c *WebsocketClient) SendMessage(message messages.WebsocketRequest) error {
	if c.closed {
		return fmt.Errorf("cannot send a message to closed websocket")
	}
	messageBytes := message.Marshal()
	slog.Debug("Sending message", "message", string(messageBytes))
	err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

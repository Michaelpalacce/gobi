package socket

import (
	"fmt"
	"io"
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
    //TODO: Implement
	closed        bool

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
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

// SendBigFile will send an item to the client without storing bigger than 1024 chunks in memory
func (c *WebsocketClient) SendItem(reader io.Reader) error {
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading: %s", err)
		}

		err = c.Conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			return fmt.Errorf("error reading file chunk: %s", err)
		}
	}

	return nil
}

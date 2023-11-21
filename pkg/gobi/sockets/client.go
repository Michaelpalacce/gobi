package sockets

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/gobi/sockets/messages"
	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
// Make sure to call Init() before using it.
type Client struct {
	Conn    *websocket.Conn
	Version int
	// Add any other fields you need for tracking the client
}

// Init will fetch information about the user. Make sure to call this first
func (c *Client) Init() {
	// Ask which version to use.
	// @TODO Partially decode the message to determine the type
	responseType, responseBody, err := c.Ask(messages.VersionType, 0, []byte(""))
	log.Println("ResponseType", responseType)
	log.Println("responseBody", string(responseBody))
	log.Println("err", err)
}

// Ask will ask a simple question to the user and wait for a response
func (c *Client) Ask(messageType string, version int, data []byte) (responseType int, responseBody []byte, err error) {
	if err = c.sendMessage(messageType, version, data); err != nil {
		return 0, nil, err
	}

	return c.Conn.ReadMessage()
}

// sendMessage enforces a uniform style in sending data
func (c *Client) sendMessage(messageType string, version int, payload interface{}) error {
	message := map[string]interface{}{
		"type":    messageType,
		"version": version,
		"payload": payload,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error encoding message: %s", err)
	}

	err = c.Conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}
	return nil
}

// Close will gracefully close the connection
func (c *Client) Close() {
	// Close the WebSocket connection gracefully
	err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		slog.Error("Error sending close message", "err", err)
	}
	time.Sleep(time.Second)

	slog.Debug("Clietn shutdown complete")
}

package sockets

import (
	"encoding/json"
	"fmt"
	"log/slog"

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
	if err := c.sendMessage(messages.NewVersionRequestPayloadMessage()); err != nil {
		slog.Error("Error while trying to get version to use", "err", err)
		return
	}

	c.readMessage()
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This message is blocking and will continue when close is called
func (c *Client) readMessage() {
out:
	for {
		messageType, message, err := c.Conn.ReadMessage()

		if err != nil {
			c.Close(fmt.Sprintf("error reading message: %", err))
			break out
		}

		switch messageType {
		case websocket.TextMessage:
			c.processTextMessage(message)
		case websocket.BinaryMessage:
		case websocket.PingMessage:
			err := c.Conn.WriteMessage(websocket.PongMessage, []byte(""))

			if err != nil {
				c.Close(fmt.Sprintf("error sending message: %s", err))
				break out
			}
		case websocket.PongMessage:
			// Do nothing, it's just a response
		case websocket.CloseMessage:
			c.Close("")
			break out
		default:
			// Do nothing, we don't know what this is.
		}
	}
}

// processTextMessage will process different types of text messages
func (c *Client) processTextMessage(message []byte) error {
	var websocketResponse messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketResponse); err != nil {
		return fmt.Errorf("error while unmarshaling websocket response %s", err)
	}

	switch websocketResponse.Version {
	case 0:
	case 1:
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketResponse.Version)
	}

	fmt.Println(websocketResponse.Payload)
	fmt.Println(websocketResponse.Type)
	fmt.Println(websocketResponse.Version)

	return nil
}

// processV0 will process all v0 websocket messages
func (c *Client) processV0(websocketResponse messages.WebsocketMessage) error {
	return nil
}

// processV1 will process all v1 websocket messages
func (c *Client) processV1(websocketResponse messages.WebsocketMessage) error {
	return nil
}

// processBinaryMessage will process different types of binary messages
func (c *Client) processBinaryMessage(message []byte) error {
	var websocketResponse messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketResponse); err != nil {
		return fmt.Errorf("error while unmarshaling websocket response %s", err)
	}

	switch websocketResponse.Version {
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketResponse.Version)
	}

	return nil
}

// sendMessage enforces a uniform style in sending data
func (c *Client) sendMessage(message messages.WebsocketMessage) error {
	err := c.Conn.WriteMessage(websocket.TextMessage, message.Marshal())

	if err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

// Close will gracefully close the connection
func (c *Client) Close(msg string) {
	payload := messages.NewCloseRequestPayloadMessage(msg)

	fmt.Print(string(payload.Marshal()))

	// Close the WebSocket connection gracefully
	err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, string(payload.Marshal())))
	if err != nil {
		slog.Error("Error sending close message", "err", err)
	}
}

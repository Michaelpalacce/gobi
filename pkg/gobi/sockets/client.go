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

// Listen will fetch initial information from the client and then listen for more data.
func (c *Client) Listen() error {
	if err := c.sendMessage(messages.NewVersionRequestPayloadMessage()); err != nil {
		slog.Error("Error while trying to get version to use", "err", err)
		return nil
	}

	return c.readMessage()
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This message is blocking and will continue when close is called
func (c *Client) readMessage() error {
	var closeError error

out:
	for {
		messageType, message, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				// The other side has closed the connection gracefully
				fmt.Println("Connection closed by the client.")
				break out
			}

			// Handle other errors
			closeError = fmt.Errorf("error reading message: %s", err)
			break out
		}

		switch messageType {
		case websocket.TextMessage:
			if closeError = c.processTextMessage(message); closeError != nil {
				break out
			}
		case websocket.BinaryMessage:
			if closeError = c.processBinaryMessage(message); closeError != nil {
				break out
			}
		case websocket.PingMessage:
			if err := c.Conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil {
				closeError = fmt.Errorf("error sending message: %s", err)
				break out
			}
		case websocket.PongMessage:
			// Do nothing, it's just a response
		case websocket.CloseMessage:
			break out
		default:
			// Do nothing, we don't know what this is.
			closeError = fmt.Errorf("error, unknown message type %d", messageType)
			break out
		}
	}

	return closeError
}

// processTextMessage will process different types of text messages
func (c *Client) processTextMessage(message []byte) error {
	var websocketResponse messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketResponse); err != nil {
		return fmt.Errorf("error while unmarshaling websocket response %s", err)
	}

	switch websocketResponse.Version {
	case 0:
		if err := c.processV0(websocketResponse); err != nil {
			return err
		}
	case 1:
		if err := c.processV1(websocketResponse); err != nil {
			return err
		}
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

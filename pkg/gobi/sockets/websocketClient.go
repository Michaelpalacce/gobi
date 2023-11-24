package sockets

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/gobi/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/gobi/messages/v1"
	"github.com/gorilla/websocket"
)

// WebsocketClient represents a connected WebSocket client
type WebsocketClient struct {
	Conn    *websocket.Conn
	Version int
	// Add any other fields you need for tracking the client
}

// Listen will request information from the client and then listen for data.
func (c *WebsocketClient) Listen() error {
	if err := c.sendMessage(messages.NewVersionRequestPayloadMessage()); err != nil {
		slog.Error("Error while trying to get version to use", "err", err)
		return nil
	}

	return c.readMessage()
}

// Close will gracefully close the connection. If an error ocurrs during closing, it will be ignored.
func (c *WebsocketClient) Close(msg string) {
	payload := messages.NewCloseRequestPayloadMessage(msg)

	fmt.Print(string(payload.Marshal()))

	// Close the WebSocket connection gracefully
	_ = c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, string(payload.Marshal())))
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This function is blocking and will stop when Close is called
func (c *WebsocketClient) readMessage() (closeError error) {
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
			if closeError = c.processPingMessage(message); closeError != nil {
				break out
			}
		case websocket.PongMessage:
			// Do nothing
		case websocket.CloseMessage:
			break out
		default:
			closeError = fmt.Errorf("error, unknown message type %d", messageType)
			break out
		}
	}

	return closeError
}

// processTextMessage will process different types of text messages
func (c *WebsocketClient) processTextMessage(message []byte) error {
	var websocketMessage messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketMessage); err != nil {
		return fmt.Errorf("error while unmarshaling websocket message %s", err)
	}

	switch websocketMessage.Version {
	case 0:

		if err := c.processV0(websocketMessage); err != nil {
			return err
		}
	case 1:
		if err := v1.ProcessMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processV0 since V0 are special, they are handled directly by the client.
// V0 messages are client specific
func (c *WebsocketClient) processV0(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	case messages.VersionType:

	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processBinaryMessage will process different types of binary messages
// TODO finish this
func (c *WebsocketClient) processBinaryMessage(message []byte) error {
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

// processPingMessage will send a PongMessage and nothing else
func (c *WebsocketClient) processPingMessage(message []byte) error {
	if err := c.Conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

// sendMessage enforces a uniform style in sending data
func (c *WebsocketClient) sendMessage(message messages.WebsocketMessage) error {
	err := c.Conn.WriteMessage(websocket.TextMessage, message.Marshal())

	if err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

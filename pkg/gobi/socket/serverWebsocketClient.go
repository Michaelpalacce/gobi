package socket

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/client"
	processor_v1 "github.com/Michaelpalacce/gobi/pkg/gobi/processor/v1"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/gorilla/websocket"
)

// ServerWebsocketClient represents a connected WebSocket client
type ServerWebsocketClient struct {
	Client *client.WebsocketClient
	// Add any other fields you need for tracking the client
}

// Listen will request information from the client and then listen for data.
func (c *ServerWebsocketClient) Listen(closeChan chan<- error) {
	closeChan <- c.readMessage()
}

// Close will gracefully close the connection. If an error ocurrs during closing, it will be ignored.
func (c *ServerWebsocketClient) Close(msg string) {
	c.Client.Close(msg)
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This function is blocking and will stop when Close is called
// NOTE: Think if we need to make each process handling async here?
func (c *ServerWebsocketClient) readMessage() (closeError error) {
out:
	for {
		messageType, message, err := c.Client.Conn.ReadMessage()
		slog.Debug("Received message from client", "message", string(message), "messageType", messageType)

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
func (c *ServerWebsocketClient) processTextMessage(message []byte) error {
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
		if err := processor_v1.ProcessServerTextMessage(websocketMessage, c.Client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processV0 since V0 are special, they are handled directly by the client.
// V0 messages are client specific
func (c *ServerWebsocketClient) processV0(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	case messages.VersionType:
		var versionResponsePayload messages.VersionPayload

		if err := json.Unmarshal(websocketMessage.Payload, &versionResponsePayload); err != nil {
			return err
		} else {
			c.Client.Client.Version = versionResponsePayload.Version
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 0", websocketMessage.Type)
	}

	return nil
}

// processBinaryMessage will process different types of binary messages
func (c *ServerWebsocketClient) processBinaryMessage(message []byte) error {
	var websocketMessage messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketMessage); err != nil {
		return fmt.Errorf("error while unmarshaling websocket response %s", err)
	}

	switch websocketMessage.Version {
	case 1:
		if err := processor_v1.ProcessServerBinaryMessage(websocketMessage, c.Client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processPingMessage will send a PongMessage and nothing else
func (c *ServerWebsocketClient) processPingMessage(message []byte) error {
	if err := c.Client.Conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

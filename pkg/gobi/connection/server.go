package connection

import (
	"encoding/json"
	"fmt"

	processor_v1 "github.com/Michaelpalacce/gobi/pkg/gobi/processor/v1"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/gorilla/websocket"
)

// ServerConnection represents a connected WebSocket client
// It will handle the processing of the messages and send them to the processor
// It will also handle the initial connection
// Reconnections are not handled here, the clients will have to handle that
type ServerConnection struct {
	WebsocketClient *socket.WebsocketClient
	V1Processor     *processor_v1.Processor
}

// Listen will request information from the client and then listen for data.
func (c *ServerConnection) Listen(closeChan chan<- error) {
	closeChan <- c.readMessage()
}

// Close will gracefully close the connection. If an error ocurrs during closing, it will be ignored.
func (c *ServerConnection) Close(msg string) {
	c.WebsocketClient.Close(msg)
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This function is blocking and will stop when Close is called
// TODO: At the same time subscribe to a Redis queue. In case something gets changed, we'll know and we'll send it over eventually.
func (c *ServerConnection) readMessage() (closeError error) {
out:
	for {
		messageType, message, err := c.WebsocketClient.Conn.ReadMessage()
		// slog.Debug("Received message from client", "message", string(message), "messageType", messageType)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				// The other side has closed the connection gracefully
				fmt.Println("Connection closed by the client.")
				break out
			}

			// Handle other errors
			closeError = fmt.Errorf("error reading message: %w", err)
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
func (c *ServerConnection) processTextMessage(message []byte) error {
	var websocketMessage messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketMessage); err != nil {
		return fmt.Errorf("error while unmarshaling websocket message %w", err)
	}

	switch websocketMessage.Version {
	case 0:
		if err := c.processV0(websocketMessage); err != nil {
			return err
		}
	case 1:
		if err := c.V1Processor.ProcessServerTextMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processV0 since V0 are special, they are handled directly by the client.
// V0 messages are client specific
func (c *ServerConnection) processV0(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	case messages.VersionType:
		var versionResponsePayload messages.VersionPayload

		if err := json.Unmarshal(websocketMessage.Payload, &versionResponsePayload); err != nil {
			return err
		} else {
			c.WebsocketClient.Client.Version = versionResponsePayload.Version
		}

		switch c.WebsocketClient.Client.Version {
		case 1:
			c.V1Processor = processor_v1.NewProcessor(c.WebsocketClient)
		default:
			return fmt.Errorf("unknown version: %d", c.WebsocketClient.Client.Version)
		}

	default:
		return fmt.Errorf("unknown websocket message type: %s for version 0", websocketMessage.Type)
	}

	return nil
}

// processBinaryMessage will process different types of binary messages
func (c *ServerConnection) processBinaryMessage(message []byte) error {
	var websocketMessage messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketMessage); err != nil {
		return fmt.Errorf("error while unmarshaling websocket response %w", err)
	}

	switch websocketMessage.Version {
	case 1:
		if err := c.V1Processor.ProcessServerBinaryMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processPingMessage will send a PongMessage and nothing else
func (c *ServerConnection) processPingMessage(message []byte) error {
	if err := c.WebsocketClient.Conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

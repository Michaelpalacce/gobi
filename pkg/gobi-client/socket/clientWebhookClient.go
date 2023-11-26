package socket

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/gorilla/websocket"
)

// ClientWebsocketClient represents a connected WebSocket client
type ClientWebsocketClient struct {
	Client *client.WebsocketClient
	// Add any other fields you need for tracking the client
}

// Listen will request information from the client and then listen for data.
func (c *ClientWebsocketClient) Listen(closeChan chan<- error) {
	initChan := make(chan error, 1)
	readMessageChan := make(chan error, 1)
	defer close(initChan)
	defer close(readMessageChan)

	go c.init(initChan)
	go c.readMessage(readMessageChan)

	select {
	case err := <-initChan:
		closeChan <- err
	case err := <-readMessageChan:
		closeChan <- err
	}
}

// Close will gracefully close the connection. If an error ocurrs during closing, it will be ignored.
func (c *ClientWebsocketClient) Close(msg string) {
	c.Client.Close(msg)
}

// init will send the initial data to the server. Stuff like what version is being used and what is the name of the vault
func (c *ClientWebsocketClient) init(initChan chan<- error) {
	if err := c.Client.SendMessage(messages.NewVersionMessage(c.Client.Client.Version)); err != nil {
		initChan <- err
		return
	}

	if err := c.Client.SendMessage(v1.NewVaultNameMessage(c.Client.Client.VaultName)); err != nil {
		initChan <- err
		return
	}

	if err := c.Client.SendMessage(v1.NewSyncMessage(c.Client.Client.LastSync)); err != nil {
		initChan <- err
		return
	}
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This function is blocking and will stop when Close is called
func (c *ClientWebsocketClient) readMessage(readMessageChan chan<- error) {
	var closeError error

out:
	for {
		messageType, message, err := c.Client.Conn.ReadMessage()
		slog.Debug("Received message from server", "message", string(message), "messageType", messageType)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				// The other side has closed the connection gracefully
				fmt.Println("Connection closed by the server.")
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

		// messageProcess := make(chan error, 1)
		//
		// go func(messageProcess chan<- error) {
		//
		// }(messageProcess)
		// // TODO: This should be in a goroutine
	}

	readMessageChan <- fmt.Errorf("error while communicating with server: %s", closeError)
}

// processTextMessage will process different types of text messages
func (c *ClientWebsocketClient) processTextMessage(message []byte) error {
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
		if err := v1.ProcessClientTextMessage(websocketMessage, c.Client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processV0 since V0 are special, they are handled directly by the client.
// V0 messages are client specific
func (c *ClientWebsocketClient) processV0(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}
}

// processBinaryMessage will process different types of binary messages
// TODO: finish this
func (c *ClientWebsocketClient) processBinaryMessage(message []byte) error {
	var websocketMessage messages.WebsocketMessage

	if err := json.Unmarshal(message, &websocketMessage); err != nil {
		return fmt.Errorf("error while unmarshaling websocket response %s", err)
	}

	switch websocketMessage.Version {
	case 1:
		if err := v1.ProcessClientBinaryMessage(websocketMessage, c.Client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processPingMessage will send a PongMessage and nothing else
func (c *ClientWebsocketClient) processPingMessage(message []byte) error {
	if err := c.Client.Conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

package connection

import (
	"encoding/json"
	"fmt"
	"log/slog"

	processor_v1 "github.com/Michaelpalacce/gobi/pkg/gobi-client/processor/v1"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/gorilla/websocket"
)

// ClientConnection handles the initial processing of the websocket messages and sends it off to the WebsocketClient to take care of them
type ClientConnection struct {
	WebsocketClient *socket.WebsocketClient
	// Add any other fields you need for tracking the client
}

// Listen will request information from the client and then listen for data.
func (c *ClientConnection) Listen(closeChan chan<- error) {
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
func (c *ClientConnection) Close(msg string) {
	c.WebsocketClient.Close(msg)
}

// init will send the initial data to the server. Stuff like what version is being used and what is the name of the vault
func (c *ClientConnection) init(initChan chan<- error) {
	if err := c.WebsocketClient.SendMessage(messages.NewVersionMessage(c.WebsocketClient.Client.Version)); err != nil {
		initChan <- err
		return
	}

	if err := c.WebsocketClient.SendMessage(v1.NewVaultNameMessage(c.WebsocketClient.Client.VaultName)); err != nil {
		initChan <- err
		return
	}

	if err := c.WebsocketClient.SendMessage(v1.NewSyncMessage(c.WebsocketClient.Client.LastSync)); err != nil {
		initChan <- err
		return
	}
}

// readMessage will continuously wait for incomming messages and process them for the given client
// This function is blocking and will stop when Close is called
func (c *ClientConnection) readMessage(readMessageChan chan<- error) {
	var closeError error

out:
	for {
		messageType, message, err := c.WebsocketClient.Conn.ReadMessage()
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
	}

	readMessageChan <- fmt.Errorf("error while communicating with server: %s", closeError)
}

// processTextMessage will process different types of text messages
func (c *ClientConnection) processTextMessage(message []byte) error {
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
		if err := processor_v1.ProcessClientTextMessage(websocketMessage, c.WebsocketClient); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket version: %d", websocketMessage.Version)
	}

	return nil
}

// processV0 since V0 are special, they are handled directly by the client.
// V0 messages are client specific
func (c *ClientConnection) processV0(websocketMessage messages.WebsocketMessage) error {
	switch websocketMessage.Type {
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}
}

// processBinaryMessage will process different types of binary messages
// When processing the binary message we need to know where to store it
func (c *ClientConnection) processBinaryMessage(message []byte) error {
	if err := processor_v1.ProcessClientBinaryMessage(message, c.WebsocketClient); err != nil {
		return err
	}

	return nil

}

// processPingMessage will send a PongMessage and nothing else
func (c *ClientConnection) processPingMessage(message []byte) error {
	if err := c.WebsocketClient.Conn.WriteMessage(websocket.PongMessage, []byte("")); err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}

	return nil
}

package socket

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/models"
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

func (c *WebsocketClient) WatchVault(changeChan chan<- *models.Item) error {
	c.StorageDriver.WatchVault(c.Client.VaultName, changeChan)
	return nil
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

// SendItem will send an item to the the client/server
func (c *WebsocketClient) SendItem(item models.Item) error {
	slog.Debug("Sending file to server", "item", item)

	reader, err := c.StorageDriver.GetReader(item)
	if err != nil {
		return err
	}
	defer reader.Close()

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

	slog.Debug("File Sent Successfully", "item", item)

	return nil
}

// FetchItem will receive an item from the client/server
func (c *WebsocketClient) FetchItem(item models.Item) error {
	slog.Debug("Fetching file", "item", item)
	c.SendMessage(v1.NewItemFetchMessage(item))

	writer, err := c.StorageDriver.GetWriter(item)
	if err != nil {
		return err
	}

	defer func() {
		writer.Close()
	}()

	if item.Size == 0 {
		slog.Debug("File Fetched Successfully", "item", item)
		return nil
	}

	bytesRead := 0
	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			return err
		}

		if messageType != websocket.BinaryMessage {
			return fmt.Errorf("invalid messageType received: %d, expected 2 (BinaryMessage)", messageType)
		}

		writer.Write(message)

		bytesRead += len(message)
		if bytesRead == item.Size {
			writer.Close()
			break
		}

		if bytesRead > item.Size {
			return fmt.Errorf("expected %d bytes, but got %d", item.Size, bytesRead)
		}
	}
	slog.Debug("File Fetched Successfully", "item", item)

	return nil
}

package processor_v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
)

// ProcessClientTextMessage will decide how to process the text message.
func ProcessClientTextMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	switch websocketMessage.Type {
	case v1.ItemSyncType:
		if err := processItemSyncMessage(websocketMessage, client); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

// processItemSyncMessage adds items to the queue
// Check if sha256 matches locally
// Request File if it does not.
func processItemSyncMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	var itemSyncPayload v1.ItemSyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemSyncPayload); err != nil {
		return err
	}

	if ok := client.StorageDriver.CheckIfLocalMatch(itemSyncPayload.Item); !ok {
		// Fetch the file
	}

	return nil
}

// ProcessClientBinaryMessage will decide how to process the binary message.
func ProcessClientBinaryMessage(message []byte, client *client.WebsocketClient) error {
	fmt.Println("Received binary message")

	// Save the received file
	err := saveFile("received_file.txt", message)
	if err != nil {
		return fmt.Errorf("error saving file: %s", err)
	}

	return nil
}

// saveFile will create the file and save the data to it
// TODO: Make it sure we don't load the entirety of the file in memory before saving it.
// TODO: what if the file already exists?
func saveFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(data))
	return err
}

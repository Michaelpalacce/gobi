package processor_v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/models"
)

// ProcessClientTextMessage will decide how to process the text message.
func ProcessClientTextMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	switch websocketMessage.Type {
	case v1.ItemsSyncType:
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
// TODO: Save filesToFetch and implement a resume mechanism for future
func processItemSyncMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	var itemsSyncPayload v1.ItemsSyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &itemsSyncPayload); err != nil {
		return err
	}

	filesToFetch := make([]models.Item, len(itemsSyncPayload.Items))

	for _, item := range itemsSyncPayload.Items {
		if ok := client.StorageDriver.CheckIfLocalMatch(item); !ok {
			filesToFetch = append(filesToFetch, item)
		}
	}

	// Foreach filesToFetch and download them all
	// TODO: Better way of handling authentication. Perhaps send it off to another class to handle
	for _, item := range filesToFetch {
		queryParams := url.Values{}
		queryParams.Add("vaultName", client.Client.VaultName)
		queryParams.Add("serverPath", item.ServerPath)

		url := url.URL{
			Scheme:   "http",
			Host:     client.Options.Host,
			Path:     "/v1/item",
			RawQuery: queryParams.Encode(),
			User:     url.UserPassword(client.Options.Username, client.Options.Password),
		}

		err := downloadFile(url, item)

		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(url url.URL, item models.Item) error {
	// Make GET request
	response, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check if the response status code is OK (200)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	// Create or truncate the file
	// TODO: Don't have me hardcoded
	file, err := os.Create(filepath.Join("./dev/clientFolder", item.ServerPath))
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("File downloaded and saved to: %s\n", item.ServerPath)
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

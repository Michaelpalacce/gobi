package processor_v1

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
)

// ProcessClientTextMessage will decide how to process the text message.
func ProcessClientTextMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
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

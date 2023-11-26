package v1

import (
	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
)

// ProcessServerBinaryMessage will decide how to process the binary message.
func ProcessServerBinaryMessage(websocketMessage messages.WebsocketMessage, client *client.WebsocketClient) error {
	return nil
}

package v1

import (
	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/messages"
)

// ProcessClientBinaryMessage will decide how to process the binary message.
func ProcessClientBinaryMessage(websocketMessage messages.WebsocketMessage, client client.WebsocketClient) error {
	return nil
}

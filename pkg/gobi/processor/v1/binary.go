package processor_v1

import (
	"fmt"

	"github.com/Michaelpalacce/gobi/pkg/messages"
)

// ProcessServerBinaryMessage will decide how to process the binary message.
func (p *Processor) ProcessServerBinaryMessage(websocketMessage messages.WebsocketMessage) error {
	return fmt.Errorf("binary messages are not supported for version 1")
}

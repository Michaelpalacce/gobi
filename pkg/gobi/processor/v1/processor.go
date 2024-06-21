package processor_v1

import (
	"github.com/Michaelpalacce/gobi/pkg/socket"
)

type Processor struct {
	WebsocketClient *socket.WebsocketClient
}

// NewProcessor will create a new processor with a default sync strategy of LastModifiedTime
// The SyncStrategy can be changed later
func NewProcessor(client *socket.WebsocketClient) *Processor {
	return &Processor{
		WebsocketClient: client,
	}
}

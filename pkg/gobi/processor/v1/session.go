package processor_v1

import (
	"github.com/Michaelpalacce/gobi/pkg/messages/v1/rest"
)

// Init will send a session id message to the client
func (p *Processor) NewSession() {
	if err := p.WebsocketClient.SendMessage(rest.NewSessionMessage(p.Session.SessionID)); err != nil {
		return
	}
}

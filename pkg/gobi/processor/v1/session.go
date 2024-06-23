package processor_v1

import (
	"github.com/Michaelpalacce/gobi/pkg/gobi/session"
	"github.com/Michaelpalacce/gobi/pkg/messages/v1/rest"
)

// Init will send a session id message to the client
func (p *Processor) NewSession() {
	// @TODO: Save the session to Redis, not sure if here, or in the `NewSession` method directly
	sessionData := session.NewSession(&p.WebsocketClient.Client, &p.WebsocketClient.User)

	if err := p.WebsocketClient.SendMessage(rest.NewSessionMessage(sessionData.SessionID)); err != nil {
		return
	}
}

package rest

import (
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
)

// ------------------------------ Session ------------------------------

type SessionPayload struct {
	SessionId string `json:"sesion_id"`
}

func NewSessionMessage(sessionId string) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: SessionType,
		Payload: SessionPayload{
			SessionId: sessionId,
		},
		Version: v1.Version,
	}
}

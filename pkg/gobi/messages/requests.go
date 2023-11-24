package messages

// CloseRequestPayload is a payload with a message that is sent to the client
type CloseRequestPayload struct {
	Reason string `json:"reason"`
}

// NewCloseRequestPayloadMessage will return a new close request message
func NewCloseRequestPayloadMessage(reason string) WebsocketMessage {
	return WebsocketMessage{
		Type: CloseRequestTyep,
		Payload: CloseRequestPayload{
			Reason: reason,
		},
		Version: 0,
	}
}

type VersionRequestPayload struct {
}

// NewClosePayloadMessage will return a new version request message
func NewVersionRequestPayloadMessage() WebsocketMessage {
	return WebsocketMessage{
		Type:    CloseRequestTyep,
		Payload: VersionRequestPayload{},
		Version: 0,
	}
}

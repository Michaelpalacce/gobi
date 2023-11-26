package messages

type CloseRequestPayload struct {
	Reason string `json:"reason"`
}

// NewCloseRequestMessage will return a new close request message
func NewCloseRequestMessage(reason string) WebsocketRequest {
	return WebsocketRequest{
		Type: CloseRequestType,
		Payload: CloseRequestPayload{
			Reason: reason,
		},
		Version: 0,
	}
}

type VersionRequestPayload struct {
}

// NewVersionRequestMessage will return a new version request message
func NewVersionRequestMessage() WebsocketRequest {
	return WebsocketRequest{
		Type:    VersionRequestType,
		Payload: VersionRequestPayload{},
		Version: 0,
	}
}

package messages

type ClosePayload struct {
	Reason string `json:"reason"`
}

// NewCloseMessage will return a new close request message
func NewCloseMessage(reason string) WebsocketRequest {
	return WebsocketRequest{
		Type: CloseType,
		Payload: ClosePayload{
			Reason: reason,
		},
		Version: 0,
	}
}

// VersionPayload is a general payload that won't change.
// It will specify what version of the websockets API to use
type VersionPayload struct {
	Version int `json:"version"`
}

func NewVersionMessage(version int) WebsocketRequest {
	return WebsocketRequest{
		Type: VersionType,
		Payload: VersionPayload{
			Version: version,
		},
		Version: 0,
	}
}

package messages

// VersionResponsePayload is a general payload that won't change.
// It will specify what version of the websockets API to use
type VersionResponsePayload struct {
	Version int `json:"version"`
}

// NewVersionResponsePayloadMessage will return a new version response message
func NewVersionResponsePayloadMessage(version int) WebsocketMessage {
	return WebsocketMessage{
		Type: VersionResponseType,
		Payload: VersionResponsePayload{
			Version: version,
		},
		Version: 0,
	}
}

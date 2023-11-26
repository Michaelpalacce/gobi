package messages

import (
	"encoding/json"
	"fmt"
)

// Holds different message types
var (
	VersionType = "version"
	CloseType   = "close"
)

// WebsocketRequest is the general WebsocketRequest that all requests will follow.
// The Payload will be the dynamic element
type WebsocketRequest struct {
	Version int         `json:"version"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Marshal will return the response as a string
func (r WebsocketRequest) Marshal() []byte {
	data, err := json.Marshal(r)

	if err != nil {
		return []byte(fmt.Errorf("could not marshal body: %s", err).Error())
	}

	return data
}

// WebsocketMessage is the general WebsocketMessage that is received by the client or server.
// The distinction from the WebsocketRequest is that the payload is RawJson and can selectively be Unmarshalled
type WebsocketMessage struct {
	Version int             `json:"version"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

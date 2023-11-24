package messages

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Holds differnt message types
var (
	VersionType = "version"

	CloseRequestTyep = "close"
)

type Client struct {
	Conn    *websocket.Conn
	Version int
}

// WebsocketMessage is the general WebsocketMessage that all requests and responses will follow.
// The Payload will be the dynamic element
type WebsocketMessage struct {
	Version int         `json:"version"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Marshal will return the response as a string
func (r WebsocketMessage) Marshal() []byte {
	data, err := json.Marshal(r)

	if err != nil {
		return []byte(fmt.Errorf("could not marshal body: %s", err).Error())
	}

	return data
}

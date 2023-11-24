package messages

// VersionPayload is a general payload that won't change.
// It will specify what version of the websockets API to use
type VersionPayload struct {
	Version int `json:"version"`
}

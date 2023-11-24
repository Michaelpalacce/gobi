package messages

// VersionResponsePayload is a general payload that won't change.
// It will specify what version of the websockets API to use
type VersionResponsePayload struct {
	Version int `json:"version"`
}

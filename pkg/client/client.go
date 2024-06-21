package client

// ClientMetadata contains metadata about the client for websocket communication
type ClientMetadata struct {
	// General
	VaultName    string
	Version      int
	LastSync     int
	SyncStrategy int
}

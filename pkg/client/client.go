package client

// ClientMetadata contains metadata about the client for websocket communication
type ClientMetadata struct {
	// General
	VaultName    string `json:"vault_name"`
	Version      int    `json:"version"`
	LastSync     int    `json:"last_sync"`
	SyncStrategy int    `json:"sync_strategy"`
}

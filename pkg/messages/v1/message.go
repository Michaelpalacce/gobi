package v1

import "github.com/Michaelpalacce/gobi/pkg/messages"

// VaultNamePayload contains the name of the vault the client wants to connect to
type VaultNamePayload struct {
	VaultName string `json:"name"`
}

// NewVaultNameMessage will return a new vault name message
func NewVaultNameMessage(vaultName string) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: VaultNameType,
		Payload: VaultNamePayload{
			VaultName: vaultName,
		},
		Version: 1,
	}
}

// SyncPayload represents a payload telling the Server that the client wants to sync.
// Contains the last time the client successfully synced with the server
type SyncPayload struct {
	// LastSync is timestamp in UTC
	LastSync int `json:"lastSync"`
}

// NewSyncMessage creates a message to send to the server telling it when was the last time the client synced.
// The LastSync is a timestamp
func NewSyncMessage(lastSync int) messages.WebsocketRequest {
	return messages.WebsocketRequest{
		Type: SyncType,
		Payload: SyncPayload{
			LastSync: lastSync,
		},
		Version: 1,
	}
}

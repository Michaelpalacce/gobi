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

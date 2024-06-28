package v1

var (
	// Client -> Server, the client tells the server which vault it wants to connect to
	VaultNameType = "vaultName"

	// Client -> Server, the client tells the server which sync strategy it wants to use
	SyncStrategyType = "syncStrategy"

	// Client -> Server, the client tells the server when was the last time it synced
	// Server -> Client, the server tells the client when was the last time it synced
	// Denotes the start of the sync process
	SyncType = "sync"

	// Server -> Client, the server tells the client what items have been modified since the last sync
	// Client -> Server, the client tells the server what items have been modified since the last sync
	SyncData = "syncData"
)

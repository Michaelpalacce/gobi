package v1

var (
	// Client -> Server, the client tells the server which vault it wants to connect to
	VaultNameType = "vaultName"

	// Client -> Server, the client tells the server which sync strategy it wants to use
	SyncStrategyType = "syncStrategy"

	// Client -> Server, Server -> Client, the client/server tells the server/client when was the last time it synced
	// Denotes the start of the sync process
	SyncType = "sync"

	// Server -> Client, Server -> Client, the server/client tells the client/server which items have changed since the last sync
	//
	InitialSyncType = "initialSync"

	// Server -> Client, Server -> Client, the server/client tells the client/server that the initial sync is done
	// Denotes the end of the sync process
	InitialSyncDoneType = "initialSyncDone"
)

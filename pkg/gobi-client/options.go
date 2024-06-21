package gobiclient

type Options struct {
	// Required
	Username         string
	Password         string
	Host             string
	VaultName        string
	VaultPath        string
	SyncStrategy     int
	WebsocketVersion int
}

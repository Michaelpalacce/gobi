package client

import "github.com/Michaelpalacce/gobi/pkg/models"

// Client contains metadata about the client
// Used by the server and client
// @TODO: Separate the client into server and client
type Client struct {
	// General
	VaultName    string
	Version      int
	LastSync     int
	SyncStrategy int

	// Server Exclusive
	User models.User
}

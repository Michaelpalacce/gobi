package client

import "github.com/Michaelpalacce/gobi/pkg/models"

// Client contains metadata about the client
type Client struct {
	User      models.User
	Version   int
	VaultName string
	VaultPath string
	LastSync  int
}

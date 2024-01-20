package v1_syncstrategies

import (
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
)

type SyncStrategy interface {
	SendSingle(client *socket.WebsocketClient, item models.Item) error
	FetchSingle(client *socket.WebsocketClient, item models.Item, conflictMode bool) error

	Fetch(client *socket.WebsocketClient) error
	FetchConflicts(client *socket.WebsocketClient) error
	FetchMultiple(client *socket.WebsocketClient, items []models.Item, conflictMode bool) error
}

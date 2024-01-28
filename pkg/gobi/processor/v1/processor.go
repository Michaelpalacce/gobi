package processor_v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	syncstrategies "github.com/Michaelpalacce/gobi/pkg/gobi/syncStrategies"
	"github.com/Michaelpalacce/gobi/pkg/messages"
	v1 "github.com/Michaelpalacce/gobi/pkg/messages/v1"
	"github.com/Michaelpalacce/gobi/pkg/redis"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage"
	syncstrats "github.com/Michaelpalacce/gobi/pkg/syncStrategies"
)

type Processor struct {
	SyncStrategy    syncstrats.SyncStrategy
	WebsocketClient *socket.WebsocketClient
}

// NewProcessor will create a new processor with a default sync strategy of LastModifiedTime
// The SyncStrategy can be changed later
func NewProcessor(client *socket.WebsocketClient) *Processor {
	return &Processor{
		WebsocketClient: client,
		SyncStrategy: syncstrategies.NewServerLastModifiedTimeSyncStrategy(
			*syncstrats.NewLastModifiedTimeSyncStrategy(client.StorageDriver, client),
		),
	}
}

// ProcessServerBinaryMessage will decide how to process the binary message.
func (p *Processor) ProcessServerBinaryMessage(websocketMessage messages.WebsocketMessage) error {
	return fmt.Errorf("binary messages are not supported for version 1")
}

// ProcessServerTextMessage will decide how to process the text message.
func (p *Processor) ProcessServerTextMessage(websocketMessage messages.WebsocketMessage) error {
	if p.WebsocketClient.Client.Version == 0 {
		return fmt.Errorf("before communications can happen, client must send %s message to specify version to use for responses", messages.VersionType)
	}

	switch websocketMessage.Type {
	// The client tells us what the vault name is
	case v1.VaultNameType:
		if err := p.processVaultNameMessage(websocketMessage); err != nil {
			return err
		}
		// The client tells us what the sync strategy is
	case v1.SyncStrategyType:
		if err := p.processSyncStrategyMessage(websocketMessage); err != nil {
			return err
		}
	case v1.SyncType:
		if err := p.processSyncMessage(websocketMessage); err != nil {
			return err
		}
	case v1.InitialSyncType:
		if err := p.processInitialSyncMessage(websocketMessage); err != nil {
			return err
		}
	case v1.InitialSyncDoneType:
		if err := p.processInitialSyncDoneMessage(websocketMessage); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown websocket message type: %s for version 1", websocketMessage.Type)
	}

	return nil
}

func (p *Processor) processSyncStrategyMessage(websocketMessage messages.WebsocketMessage) error {
	var syncStrategyPayload v1.SyncStrategyPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncStrategyPayload); err != nil {
		return err
	}

	switch syncStrategyPayload.SyncStrategy {
	case syncstrats.LastModifiedTimeStrategy:
		p.WebsocketClient.Client.SyncStrategy = syncStrategyPayload.SyncStrategy
		p.SyncStrategy = syncstrategies.NewServerLastModifiedTimeSyncStrategy(
			*syncstrats.NewLastModifiedTimeSyncStrategy(p.WebsocketClient.StorageDriver, p.WebsocketClient),
		)
	default:
		return fmt.Errorf("unknown sync strategy: %d", syncStrategyPayload.SyncStrategy)
	}

	return nil
}

// processVaultNameMessage will set the VaultName in the client if when it's sent
func (p *Processor) processVaultNameMessage(websocketMessage messages.WebsocketMessage) error {
	var vaultNamePayload v1.VaultNamePayload

	if err := json.Unmarshal(websocketMessage.Payload, &vaultNamePayload); err != nil {
		return err
	}

	p.WebsocketClient.Client.VaultName = vaultNamePayload.VaultName

	return nil
}

// processInitialSyncMessage adds items to the queue
// This is done only once, after the initial sync, the client will watch the vault for changes
func (p *Processor) processInitialSyncMessage(websocketMessage messages.WebsocketMessage) error {
	var initialSyncPayload v1.InitialSyncPayload

	p.WebsocketClient.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer p.WebsocketClient.Conn.SetReadDeadline(time.Time{})

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncPayload); err != nil {
		return err
	}

	p.WebsocketClient.SendMessage(v1.NewInitialSyncDoneMessage(p.WebsocketClient.Client.LastSync))

	slog.Info("Fully synced")

	go p.subscribeToRedis()

	return nil
}

func (p *Processor) subscribeToRedis() {
	chanName := p.WebsocketClient.Client.User.Username + "-" + p.WebsocketClient.Client.VaultName
	redisChan := redis.Subscribe(chanName).Channel()
	slog.Info("Subscribed to Redis channel", "channel", chanName)

	for {
		msg := <-redisChan
		fmt.Println(msg.Payload)
	}
}

func (p *Processor) processInitialSyncDoneMessage(websocketMessage messages.WebsocketMessage) error {
	var initialSyncDonePayload v1.InitialSyncDonePayload

	if err := json.Unmarshal(websocketMessage.Payload, &initialSyncDonePayload); err != nil {
		return err
	}

	p.WebsocketClient.InitialSync = true
	// This is just for info
	p.WebsocketClient.Client.LastSync = initialSyncDonePayload.LastSync

	p.WebsocketClient.SendMessage(v1.NewSyncMessage(initialSyncDonePayload.LastSync))

	slog.Info("Initial Client Sync Done", "vaultName", p.WebsocketClient.Client.VaultName)

	return nil
}

// processSyncMessage will enqueue items since the last sync and send the metadata to the client
func (p *Processor) processSyncMessage(websocketMessage messages.WebsocketMessage) error {
	var syncPayload v1.SyncPayload

	if err := json.Unmarshal(websocketMessage.Payload, &syncPayload); err != nil {
		return err
	}

	p.WebsocketClient.StorageDriver.EnqueueItemsSince(
		syncPayload.LastSync,
		p.WebsocketClient.Client.VaultName,
	)

	items := p.WebsocketClient.StorageDriver.GetAllItems(storage.ConflictModeNo)

	slog.Debug("Items found for sync since last reconcillation", "items", items, "lastSync", syncPayload.LastSync)

	p.WebsocketClient.SendMessage(v1.NewInitialSyncMessage(items))

	return nil
}

package syncstrategies

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/redis"
	syncstrategies "github.com/Michaelpalacce/gobi/pkg/syncStrategies"
)

// ServerLastModifiedTimeSyncStrategy is a wrapper around the syncStrategies.LastModifiedTimeSyncStrategy
// that introduces locking with redis to prevent multiple instances of the application from
// attempting to sync the same data at the same time.
type ServerLastModifiedTimeSyncStrategy struct {
	strategy syncstrategies.LastModifiedTimeSyncStrategy
}

func NewServerLastModifiedTimeSyncStrategy(strategy syncstrategies.LastModifiedTimeSyncStrategy) *ServerLastModifiedTimeSyncStrategy {
	return &ServerLastModifiedTimeSyncStrategy{
		strategy: strategy,
	}
}

func (s *ServerLastModifiedTimeSyncStrategy) getLockKey() string {
	return fmt.Sprintf(
		"gobi:sync:%s",
		s.strategy.Client.Client.User.Username+s.strategy.Client.Client.VaultName,
	)
}

// lock will lock the sync strategy
// This is used to ensure that only one sync strategy is SendSingle or FetchSingle at a time
func (s *ServerLastModifiedTimeSyncStrategy) lock() error {
	lockKey := s.getLockKey()
	start := time.Now()
	for {
		if time.Since(start) > 30*time.Minute {
			return fmt.Errorf("waited too long to receive lock %s", lockKey)
		}

		gotLock, err := redis.Lock(lockKey, 60*time.Minute)
		if err != nil {
			return err
		}

		if gotLock {
			return nil
		}

		slog.Debug("Waiting for lock", "lockKey", lockKey)
		time.Sleep(1 * time.Second)
	}
}

// unlock will unlock the sync strategy
// Will attempt to unlock 3 times before returning an error
func (s *ServerLastModifiedTimeSyncStrategy) unlock() error {
	slog.Debug("Unlocking", "lockKey", s.getLockKey())
	for i := 0; i < 3; i++ {
		err := redis.Unlock(s.getLockKey())
		if err == nil {
			return nil
		}
		slog.Warn("Failed to unlock", "lockKey", s.getLockKey(), "err", err)
	}

	return fmt.Errorf("failed to unlock %s", s.getLockKey())
}

func (s *ServerLastModifiedTimeSyncStrategy) SendSingle(item models.Item) error {
	s.lock()
	defer s.unlock()
	return s.strategy.SendSingle(item)
}

func (s *ServerLastModifiedTimeSyncStrategy) Fetch() error {
	s.lock()
	defer s.unlock()
	return s.strategy.Fetch()
}

func (s *ServerLastModifiedTimeSyncStrategy) FetchConflicts() error {
	s.lock()
	defer s.unlock()
	return s.strategy.FetchConflicts()
}

func (s *ServerLastModifiedTimeSyncStrategy) FetchSingle(item models.Item, conflictMode bool) error {
	s.lock()
	defer s.unlock()
	return s.strategy.FetchSingle(item, conflictMode)
}

func (s *ServerLastModifiedTimeSyncStrategy) FetchMultiple(items []models.Item, conflictMode bool) error {
	s.lock()
	defer s.unlock()
	return s.strategy.FetchMultiple(items, conflictMode)
}

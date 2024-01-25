package syncstrategies

import (
	"sync"

	"github.com/Michaelpalacce/gobi/pkg/models"
	syncstrategies "github.com/Michaelpalacce/gobi/pkg/syncStrategies"
)

// ClientLastModifiedTimeSyncStrategy is a wrapper around the syncStrategies.LastModifiedTimeSyncStrategy
// that introduces locking with mutex
type ClientLastModifiedTimeSyncStrategy struct {
	strategy syncstrategies.LastModifiedTimeSyncStrategy
	mutex    sync.Mutex
}

func NewClientLastModifiedTimeSyncStrategy(strategy syncstrategies.LastModifiedTimeSyncStrategy) *ClientLastModifiedTimeSyncStrategy {
	return &ClientLastModifiedTimeSyncStrategy{
		strategy: strategy,
		mutex:    sync.Mutex{},
	}
}

// lock will lock the sync strategy
// This is used to ensure that only one sync strategy is SendSingle or FetchSingle at a time
func (s *ClientLastModifiedTimeSyncStrategy) lock() {
	s.mutex.Lock()
}

// unlock will unlock the sync strategy
func (s *ClientLastModifiedTimeSyncStrategy) unlock() {
	s.mutex.Unlock()
}

func (s *ClientLastModifiedTimeSyncStrategy) SendSingle(item models.Item) error {
	s.lock()
	defer s.unlock()
	return s.strategy.SendSingle(item)
}

func (s *ClientLastModifiedTimeSyncStrategy) Fetch() error {
	s.lock()
	defer s.unlock()
	return s.strategy.Fetch()
}

func (s *ClientLastModifiedTimeSyncStrategy) FetchConflicts() error {
	s.lock()
	defer s.unlock()
	return s.strategy.FetchConflicts()
}

func (s *ClientLastModifiedTimeSyncStrategy) FetchSingle(item models.Item, conflictMode bool) error {
	s.lock()
	defer s.unlock()
	return s.strategy.FetchSingle(item, conflictMode)
}

func (s *ClientLastModifiedTimeSyncStrategy) FetchMultiple(items []models.Item, conflictMode bool) error {
	s.lock()
	defer s.unlock()
	return s.strategy.FetchMultiple(items, conflictMode)
}

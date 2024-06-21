package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

type SyncData struct {
	// Sync Relevant Data
	LastSync int `json:"lastSync,omitempty"`
}

// readSyncData reads and then returns the sync data from the given path

func readSyncData(path string) (*SyncData, error) {
	syncBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading sync file: %w", err)
	}

	sync := &SyncData{}

	err = json.Unmarshal(syncBytes, sync)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling sync: %w", err)
	}

	return sync, nil
}

// writeSyncData writes the given sync data to the given path
// 640 so that only the owner can read and write, group can read, and others can't do anything
func writeSyncData(path string, sync *SyncData) error {
	syncBytes, err := json.Marshal(sync)
	if err != nil {
		return fmt.Errorf("error marshalling sync: %w", err)
	}

	fmt.Println(string(syncBytes))

	err = os.WriteFile(path, syncBytes, 0o640)
	if err != nil {
		return fmt.Errorf("error writing sync file: %w", err)
	}

	return nil
}

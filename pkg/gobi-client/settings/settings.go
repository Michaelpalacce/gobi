package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

type SettingsData struct {
	// The settings of the client
	VaultName        string `json:"vaultName,omitempty"`
	WebsocketVersion int    `json:"websocketVersion,omitempty"`
	SyncStrategy     int    `json:"syncStrategy,omitempty"`
}

// readSettings reads and then returns the settings from the given path
func readSettings(path string) (*SettingsData, error) {
	settingsBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading settings file: %w", err)
	}

	settings := &SettingsData{}
	err = json.Unmarshal(settingsBytes, settings)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling settings: %w", err)
	}

	return settings, nil
}

// writeSettings writes the given settings to the given path
// 640 so that only the owner can read and write, group can read, and others can't do anything
func writeSettings(path string, settings *SettingsData) error {
	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("error marshalling settings: %w", err)
	}

	err = os.WriteFile(path, settingsBytes, 0o640)
	if err != nil {
		return fmt.Errorf("error writing settings file: %w", err)
	}

	return nil
}

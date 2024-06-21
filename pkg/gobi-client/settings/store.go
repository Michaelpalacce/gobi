package settings

import (
	"fmt"
	"os"

	gobiclient "github.com/Michaelpalacce/gobi/pkg/gobi-client"
)

type Store struct {
	VaultPath string
	options   gobiclient.Options

	Settings *SettingsData
	Sync     *SyncData
}

// NewStore creates a new Store
// It takes the options as an argument. The options will be used and persisted
func NewStore(options gobiclient.Options) (*Store, error) {
	localStore := &Store{
		options: options,

		VaultPath: fmt.Sprintf("%s/%s", options.VaultPath, options.VaultName),
		Settings:  &SettingsData{},
		Sync:      &SyncData{},
	}

	err := localStore.Init()
	if err != nil {
		return nil, fmt.Errorf("error initializing LocalStore: %s", err)
	}

	return localStore, nil
}

// Init function will initialize the LocalStore
// It is responsible for creating the settings file if it doesn't exist
func (l *Store) Init() error {
	// Create the settings file if it doesn't exist
	_, err := os.Stat(l.getConfigDir())
	if os.IsNotExist(err) {
		err = os.Mkdir(l.getConfigDir(), 0o700)
		if err != nil {
			return fmt.Errorf("error creating config dir: %w", err)
		}
	}

	_, err = os.Stat(l.GetSettingsPath())
	if os.IsNotExist(err) {
		if l.options.WebsocketVersion == 0 {
			l.options.WebsocketVersion = 1
		}

		if l.options.SyncStrategy == 0 {
			l.options.SyncStrategy = 1
		}

		l.Settings.WebsocketVersion = l.options.WebsocketVersion
		l.Settings.VaultName = l.options.VaultName
		l.Settings.SyncStrategy = l.options.SyncStrategy

		err = writeSettings(l.GetSettingsPath(), l.Settings)
		if err != nil {
			return fmt.Errorf("error writing settings: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking settings file: %w", err)
	}

	_, err = os.Stat(l.GetSyncPath())
	if os.IsNotExist(err) {
		l.Sync.LastSync = 0

		err = writeSyncData(l.GetSyncPath(), l.Sync)
		if err != nil {
			return fmt.Errorf("error writing sync: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking sync file: %w", err)
	}

	settings, err := readSettings(l.GetSettingsPath())
	if err != nil {
		return fmt.Errorf("error reading settings file: %w", err)
	}
	l.Settings = settings

	sync, err := readSyncData(l.GetSyncPath())
	if err != nil {
		return fmt.Errorf("error reading sync file: %w", err)
	}
	l.Sync = sync

	return nil
}

// getConfigDir returns the path where the hidden config dir is located
// The config dir is used to store settings and other configuration files
func (l *Store) getConfigDir() string {
	return fmt.Sprintf("%s/.gobi", l.VaultPath)
}

// GetSettingsPath returns the path to the settings file
// The settings file is used to store the settings of the client
func (l *Store) GetSettingsPath() string {
	return fmt.Sprintf("%s/settings.json", l.getConfigDir())
}

// GetSyncPath returns the path to the sync file
// The sync file is used to store the last sync time of the client
func (l *Store) GetSyncPath() string {
	return fmt.Sprintf("%s/sync.json", l.getConfigDir())
}

func (l *Store) SaveSettings() error {
	return writeSettings(l.GetSettingsPath(), l.Settings)
}

func (l *Store) SaveSync() error {
	return writeSyncData(l.GetSyncPath(), l.Sync)
}

// SaveAll saves all the data in the LocalStore
func (l *Store) SaveAll() error {
	err := l.SaveSettings()
	if err != nil {
		return fmt.Errorf("error saving settings data: %w", err)
	}

	err = l.SaveSync()
	if err != nil {
		return fmt.Errorf("error saving sync data: %w", err)
	}

	return nil
}

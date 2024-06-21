package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/client"
	gobiclient "github.com/Michaelpalacce/gobi/pkg/gobi-client"
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/auth"
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/connection"
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/settings"
	"github.com/Michaelpalacce/gobi/pkg/logger"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage"
	"github.com/gorilla/websocket"
)

func main() {
	logger.ConfigureLogging()

	// Define command-line flags for username and password
	var (
		username     string
		host         string
		password     string
		vaultName    string
		vaultPath    string
		syncStrategy int
		gobiClient   *connection.ClientConnection
	)

	flag.StringVar(&host, "host", "localhost:8080", "Target host")
	flag.StringVar(&username, "username", "root", "Username for authentication")
	flag.StringVar(&password, "password", "toor", "Password for authentication")
	flag.StringVar(&vaultName, "vaultName", "testVault", "The name of the vault to connect to")
	flag.StringVar(&vaultPath, "vaultPath", ".dev/clientFolder", "The path to the vault to watch")
	flag.IntVar(&syncStrategy, "syncStrategy", 1, "The sync strategy to use. Available: 1 (default): lastModified")

	// Parse command-line flags
	flag.Parse()

	go interrupt(gobiClient)

out:
	for {
		closeChan := make(chan error, 1)
		defer close(closeChan)

		var (
			conn *websocket.Conn
			err  error
		)

		options := gobiclient.Options{
			Username:  username,
			Password:  password,
			Host:      host,
			VaultName: vaultName,
			VaultPath: vaultPath,
			// This is intenionally hardcoded to 1
			// We want to always use the latest :)
			WebsocketVersion: 1,
		}

		if conn, err = establishConn(options); err != nil {
			slog.Error("Error while trying to establish connection to server. Since no initial connection could be established, closing.", "error", err)
			break out
		}

		settingsStore, err := settings.NewStore(options)
		if err != nil {
			slog.Error("Error creating settings store", "error", err)
			break out
		}

		gobiClient = &connection.ClientConnection{
			LocalSettings: settingsStore,
			WebsocketClient: &socket.WebsocketClient{
				Client: client.ClientMetadata{
					Version:      settingsStore.Settings.WebsocketVersion,
					VaultName:    settingsStore.Settings.VaultName,
					LastSync:     settingsStore.Sync.LastSync,
					SyncStrategy: settingsStore.Settings.SyncStrategy,
				},
				Conn:          conn,
				StorageDriver: storage.NewLocalDriver(options.VaultPath),
				User: models.User{
					Username: options.Username,
					Password: options.Password,
				},
			},
		}

		go gobiClient.Listen(closeChan)

		err = <-closeChan
		if err != nil {
			slog.Error("Closing connection due to error with server", "error", err)
			gobiClient.Close(err.Error())
		}

		gobiClient.Close("")
		time.Sleep(5 * time.Second)
	}
}

// establishConn establish a connection to the server using the given option.
// Supprts only BasicAuth
func establishConn(options gobiclient.Options) (*websocket.Conn, error) {
	url := url.URL{
		Scheme: "ws",
		Host:   options.Host,
		Path:   fmt.Sprintf("/api/v%d/ws/", options.WebsocketVersion),
	}

	header := http.Header{"Authorization": []string{auth.BasicAuth(options.Username, options.Password)}}
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}

	// Establish a WebSocket connection with headers
	conn, resp, websocketErr := dialer.Dial(url.String(), header)
	if websocketErr != nil {
		if resp != nil {
			// Read and print the body content
			body, err := io.ReadAll(io.Reader(resp.Body))
			if err != nil {
				return nil, fmt.Errorf("error connecting to WebSocket: %w", websocketErr)
			}
			return nil, fmt.Errorf("error connecting to WebSocket: %w, response was %s", websocketErr, body)
		}

		return nil, fmt.Errorf("error connecting to WebSocket: %w", websocketErr)
	}

	return conn, nil
}

func interrupt(gobiClient *connection.ClientConnection) {
	// Set up a channel to handle signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

	<-interrupt

	if gobiClient != nil {
		gobiClient.Close("os.Interrupt received. Closing connection.")
	}

	os.Exit(1)
}

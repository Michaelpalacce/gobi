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
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/auth"
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/connection"
	"github.com/Michaelpalacce/gobi/pkg/logger"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/Michaelpalacce/gobi/pkg/storage"
	syncstrategies "github.com/Michaelpalacce/gobi/pkg/syncStrategies"
	"github.com/gorilla/websocket"
)

func main() {
	logger.ConfigureLogging()

	// Define command-line flags for username and password
	var (
		username   string
		host       string
		password   string
		vaultName  string
		vaultPath  string
		version    int
		gobiClient *connection.ClientConnection
	)

	flag.StringVar(&host, "host", "localhost:8080", "Target host")
	flag.StringVar(&username, "username", "test", "Username for authentication")
	flag.StringVar(&password, "password", "test", "Password for authentication")
	flag.StringVar(&vaultName, "vaultName", "testVault", "The name of the vault to connect to")
	flag.StringVar(&vaultPath, "vaultPath", ".dev/clientFolder", "The path to the vault to watch")
	flag.IntVar(&version, "version", 1, "The version of the vault to watch")

	// Parse command-line flags
	flag.Parse()

	// Check if both username and password are provided
	// TODO: Make me better... maybe a struct? or a map? or something else? I don't know yet :( I'm just a simple if statement
	if username == "" || password == "" || host == "" || vaultName == "" || vaultPath == "" || version == 0 {
		flag.Usage()
		os.Exit(1)
	}

	go interrupt(gobiClient)

out:
	for {
		closeChan := make(chan error, 1)
		defer close(closeChan)

		var (
			conn *websocket.Conn
			err  error
		)

		options := connection.Options{
			Username:  username,
			Password:  password,
			Host:      host,
			VaultName: vaultName,
			VaultPath: vaultPath,
			Version:   version,
		}

		if conn, err = establishConn(options); err != nil {
			slog.Error("Error while trying to establish connection to server. Since no initial connection could be established, closing.", "error", err)
			break out
		}

		gobiClient = &connection.ClientConnection{
			WebsocketClient: &socket.WebsocketClient{
				// TODO: Make this configurable
				// TODO: Fetch me from somewhere... maybe a file? maybe a database? maybe a configmap? maybe a secret? maybe a flag? maybe an env var?
				Client: client.Client{
					// Intentionally hardcoded to latest.
					Version:      options.Version,
					VaultName:    options.VaultName,
					LastSync:     0,
					SyncStrategy: syncstrategies.LastModifiedTimeStrategy,
					User: models.User{
						Username: options.Username,
						Password: options.Password,
					},
				},
				Conn:          conn,
				StorageDriver: storage.NewLocalDriver(options.VaultPath),
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
func establishConn(options connection.Options) (*websocket.Conn, error) {
	url := url.URL{
		Scheme: "ws",
		Host:   options.Host,
		Path:   fmt.Sprintf("/api/v%d/ws/", options.Version),
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

package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"io"

	"github.com/Michaelpalacce/gobi/internal/gobi-client/connection"
	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/socket"
	"github.com/Michaelpalacce/gobi/pkg/logger"
	"github.com/gorilla/websocket"
)

func main() {
	logger.ConfigureLogging()

	// Define command-line flags for username and password
	var (
		username  string
		host      string
		password  string
		vaultName string
		vaultPath string
	)

	flag.StringVar(&host, "host", "localhost:8080", "Target host")
	flag.StringVar(&username, "username", "test", "Username for authentication")
	flag.StringVar(&password, "password", "test", "Password for authentication")
	flag.StringVar(&vaultName, "vaultName", "testVault", "The name of the vault to connect to")
	flag.StringVar(&vaultPath, "vaultPath", "./test", "The path to the vault to watch")

	// Parse command-line flags
	flag.Parse()

	// Check if both username and password are provided
	// TODO: Make me better
	if username == "" || password == "" || host == "" || vaultName == "" || vaultPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Set up a channel to handle signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer close(interrupt)

out:
	for {
		closeChan := make(chan error, 1)
		defer close(closeChan)

		var (
			conn *websocket.Conn
			err  error
		)

		options := connection.Options{
			Username: username,
			Password: password,
			Host:     host,
		}

		if conn, err = establishConn(options); err != nil {
			slog.Error("Error while trying to establish connection to server", "error", err)
			time.Sleep(5 * time.Second)

			continue out
		}

		// establish connection to the server
		client := socket.ClientWebhookClient{
			Client: &client.WebsocketClient{
				Client: client.Client{
					// Intentionally hardcoded to latest.
					Version:   1,
					VaultName: vaultName,
					VaultPath: vaultPath,
				},
				Conn: conn,
			},
		}

		go client.Listen(closeChan)

		select {
		case err := <-closeChan:
			if err != nil {
				slog.Error("Closing connection due to error with server", "error", err)
				client.Close(err.Error())
			}

			client.Close("")
		case <-interrupt:
			client.Close("os.Interrupt received. Closing connection.")
			break out
		}

		time.Sleep(5 * time.Second)
	}

}

// establishConn will fetch the username and password from the arguments and establish a connection to the server
func establishConn(options connection.Options) (*websocket.Conn, error) {
	url := url.URL{Scheme: "ws", Host: options.Host, Path: "/ws/"}
	header := http.Header{"Authorization": []string{basicAuth(options.Username, options.Password)}}
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
				return nil, fmt.Errorf("error connecting to WebSocket: %s", websocketErr)
			}
			return nil, fmt.Errorf("error connecting to WebSocket: %s, response was %s", websocketErr, body)
		}

		return nil, fmt.Errorf("error connecting to WebSocket: %s", websocketErr)
	}

	return conn, nil
}

// basicAuth returns the Basic Authentication string
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

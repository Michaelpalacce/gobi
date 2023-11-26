package main

import (
	"encoding/base64"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"io"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/gobi-client/socket"
	"github.com/Michaelpalacce/gobi/pkg/logger"
	"github.com/gorilla/websocket"
)

func main() {
	logger.ConfigureLogging()
	// establish connection to the server
	//TODO: Client props should be loaded
	client := socket.ClientWebhookClient{
		Client: &client.WebsocketClient{
			Client: client.Client{
				Version:   1,
				VaultName: "Test",
			},
			Conn: establishConn(),
		},
	}

	// Set up a channel to handle signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	closeChan := make(chan error, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start a goroutine to read messages from the WebSocket connection
	go client.Listen(closeChan)
	defer close(closeChan)
	defer close(interrupt)

	select {
	case err := <-closeChan:
		if err != nil {
			slog.Error("Closing connection due to error with server", "error", err)
			client.Close(err.Error())
		}

		client.Close("")
	case <-interrupt:
		client.Close("os.Interrupt received. Closing connection.")
	}
}

// establishConn will fetch the username and password from the arguments and establish a connection to the server
func establishConn() *websocket.Conn {
	// Define command-line flags for username and password
	var (
		username string
		host     string
		password string
	)

	flag.StringVar(&host, "host", "localhost:8080", "Target host")
	flag.StringVar(&username, "username", "", "Username for authentication")
	flag.StringVar(&password, "password", "", "Password for authentication")

	// Parse command-line flags
	flag.Parse()

	// Check if both username and password are provided
	if username == "" || password == "" || host == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create a WebSocket connection URL
	url := url.URL{Scheme: "ws", Host: host, Path: "/ws/"}
	// Set up basic authentication
	header := http.Header{"Authorization": []string{basicAuth(username, password)}}
	// Create a WebSocket dialer
	dialer := websocket.DefaultDialer

	// Establish a WebSocket connection with headers
	conn, resp, websocketErr := dialer.Dial(url.String(), header)
	if websocketErr != nil {
		// Read and print the body content
		body, err := io.ReadAll(io.Reader(resp.Body))
		if err != nil {
			log.Fatalf("Error connecting to WebSocket: %s", websocketErr)
		} else {
			log.Fatalf("Error connecting to WebSocket: %s, response was %s", websocketErr, body)
		}
	}

	return conn
}

// basicAuth returns the Basic Authentication string
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"io"

	"github.com/Michaelpalacce/gobi/pkg/gobi-client/sockets"
	"github.com/Michaelpalacce/gobi/pkg/logger"
	"github.com/gorilla/websocket"
)

func main() {
	logger.ConfigureLogging()
	// establish connection to the server
	client := sockets.WebsocketClient{
		Version: 1, // TODO Make me dynamic
		Conn:    establishConn(),
	}

	// Set up a channel to handle signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start a goroutine to read messages from the WebSocket connection
	go client.Listen()

	// Wait for a signal to gracefully close the connection
	<-interrupt
	client.Close("")
}

// establishConn will fetch the username and password from the arguments and establish a connection to the server
// TODO: Move auth to the gobi-client
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
	if username == "" || password == "" {
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

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"io"

	"github.com/gorilla/websocket"
)

func main() {
	// establish connection to the server
	conn := establishConn()
	defer conn.Close()

	// Set up a channel to handle signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start a goroutine to read messages from the WebSocket connection
	go readMessages(conn)

	// Wait for a signal to gracefully close the connection
	<-interrupt
	shutdown(conn)
}

// establishConn will fetch the username and password from the arguments and establish a connection to the server
func establishConn() *websocket.Conn {
	// Define command-line flags for username and password
	var (
		username string
		host     string
		password string
	)

	flag.StringVar(&host, "host", "", "Target host")
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
	url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/"}
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

// shutdown will attempt to gracefully shut down a connection when the client is killed
func shutdown(conn *websocket.Conn) {
	log.Println("Interrupt signal received, closing WebSocket connection...")
	// Close the WebSocket connection gracefully
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Error sending close message:", err)
	}
	time.Sleep(time.Second)
	log.Println("Client shutdown complete.")
}

// readMessages is an infinite loop to listen for messages from the server
func readMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		conn.WriteMessage(websocket.TextMessage, json.RawMessage(`{"version":0,"type":"version","payload":{"version":1}}`))
		fmt.Printf("Received message: %s\n", message)
	}
}

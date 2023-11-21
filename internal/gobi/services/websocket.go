package services

import (
	"log"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

// connectedClientsMutex is a mutex so we can register only one client at a time
var connectedClientsMutex sync.Mutex

// Client represents a connected WebSocket client
type Client struct {
	conn *websocket.Conn
	// Add any other fields you need for tracking the client
}

// WebsocketService handles the connection between the server and the clients
type WebsocketService struct {
	// connectedClients is a map of all the connected clients
	connectedClients map[*Client]bool
}

func NewWebsocketService() WebsocketService {
	return WebsocketService{
		connectedClients: make(map[*Client]bool),
	}
}

func (s *WebsocketService) HandleConnection(conn *websocket.Conn) {
	client := &Client{conn: conn}

	s.registerClient(client)

	// Handle incoming messages
	go s.readMessages(client)

	// Example: Send a welcome message to the connected client
	err := conn.WriteMessage(websocket.TextMessage, []byte("OK"))

	if err != nil {
		slog.Error("Error sending welcome message", "error", err)
		return
	}
}

// readMessages reads incoming messages from the client
func (s *WebsocketService) readMessages(client *Client) {
	for {
		_, _, err := client.conn.ReadMessage()
		// messageType, p, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// // Handle the incoming message based on the messageType and content (p)
		// handleMessage(client, messageType, p)
	}
}

// registerClient registers a client
func (s *WebsocketService) registerClient(client *Client) {
	connectedClientsMutex.Lock()

	defer connectedClientsMutex.Unlock()
	s.connectedClients[client] = true
}

// unregisterClient unregisters a client
func (s *WebsocketService) unregisterClient(client *Client) {
	connectedClientsMutex.Lock()
	defer connectedClientsMutex.Unlock()
	delete(s.connectedClients, client)
}

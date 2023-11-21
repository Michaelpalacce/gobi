package services

import (
	"log/slog"
	"sync"

	"github.com/Michaelpalacce/gobi/pkg/gobi/sockets"
	"github.com/gorilla/websocket"
)

// connectedClientsMutex is a mutex so we can register only one client at a time
var connectedClientsMutex sync.Mutex

// WebsocketService handles the connection between the server and the clients
type WebsocketService struct {
	// connectedClients is a map of all the connected clients
	connectedClients map[*sockets.Client]bool
}

func NewWebsocketService() WebsocketService {
	return WebsocketService{
		connectedClients: make(map[*sockets.Client]bool),
	}
}

// HandleConnection will register a new client and send a Welcome Message
func (s *WebsocketService) HandleConnection(conn *websocket.Conn) {
	client := &sockets.Client{Conn: conn}

	s.registerClient(client)

	defer s.unregisterClient(client)

	s.readMessages(client)
}

func (s *WebsocketService) readMessages(client *sockets.Client) {
	for {
		_, _, err := client.Conn.ReadMessage()
		// messageType, p, err := client.conn.ReadMessage()
		if err != nil {
			slog.Error("error reading message", "err", err)
			break
		}

		// // Handle the incoming message based on the messageType and content (p)
		// handleMessage(client, messageType, p)
	}
}

// registerClient registers a client
func (s *WebsocketService) registerClient(client *sockets.Client) {
	connectedClientsMutex.Lock()

	client.Init()

	defer connectedClientsMutex.Unlock()
	s.connectedClients[client] = true
}

// unregisterClient unregisters a client
func (s *WebsocketService) unregisterClient(client *sockets.Client) {
	connectedClientsMutex.Lock()

	client.Close()

	defer connectedClientsMutex.Unlock()
	delete(s.connectedClients, client)
}

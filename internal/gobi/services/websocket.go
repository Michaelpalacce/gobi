package services

import (
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

// registerClient registers a client
func (s *WebsocketService) registerClient(client *sockets.Client) {
	connectedClientsMutex.Lock()

	s.connectedClients[client] = true

	defer connectedClientsMutex.Unlock()
}

// unregisterClient unregisters a client
func (s *WebsocketService) unregisterClient(client *sockets.Client) {
	connectedClientsMutex.Lock()

	delete(s.connectedClients, client)

	defer connectedClientsMutex.Unlock()
}

// HandleConnection will register a new client and send a Welcome Message
func (s *WebsocketService) HandleConnection(conn *websocket.Conn) {
	client := &sockets.Client{Conn: conn}

	s.registerClient(client)
	defer s.unregisterClient(client)

	err := client.Listen()

	if err != nil {
		client.Close(err.Error())
	}

	client.Close("")
}

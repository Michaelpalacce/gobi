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
	connectedClients map[*sockets.WebsocketClient]bool
}

func NewWebsocketService() WebsocketService {
	return WebsocketService{
		connectedClients: make(map[*sockets.WebsocketClient]bool),
	}
}

// HandleConnection will register a new client and start listening for any messages
// At the end, the client will be unregistered and the connection will be closed with
// an Error message if one was present
func (s *WebsocketService) HandleConnection(conn *websocket.Conn) {
	client := &sockets.WebsocketClient{Conn: conn}

	s.registerClient(client)
	defer s.unregisterClient(client)

	err := client.Listen()

	if err != nil {
		client.Close(err.Error())
	}

	client.Close("")
}

// registerClient registers a client
func (s *WebsocketService) registerClient(client *sockets.WebsocketClient) {
	connectedClientsMutex.Lock()

	s.connectedClients[client] = true

	defer connectedClientsMutex.Unlock()
}

// unregisterClient unregisters a client
func (s *WebsocketService) unregisterClient(client *sockets.WebsocketClient) {
	connectedClientsMutex.Lock()

	delete(s.connectedClients, client)

	defer connectedClientsMutex.Unlock()
}

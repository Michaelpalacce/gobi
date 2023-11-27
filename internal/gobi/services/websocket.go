package services

import (
	"log/slog"
	"sync"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/gobi/connection"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/socket"
	"github.com/gorilla/websocket"
)

// connectedClientsMutex is a mutex so we can register only one client at a time
var connectedClientsMutex sync.Mutex

// WebsocketService handles the connection between the server and the clients
type WebsocketService struct {
	// connectedClients is a map of all the connected clients
	connectedClients map[*connection.ServerConnection]bool

	DB *database.Database
}

func NewWebsocketService(db *database.Database) WebsocketService {
	return WebsocketService{
		connectedClients: make(map[*connection.ServerConnection]bool),
		DB:               db,
	}
}

// HandleConnection will register a new client and start listening for any messages
// At the end, the client will be unregistered and the connection will be closed with
// an Error message if one was present
func (s *WebsocketService) HandleConnection(conn *websocket.Conn, user models.User) {
	client := &connection.ServerConnection{Client: &socket.WebsocketClient{
		DB:   s.DB,
		Conn: conn,

		Client: client.Client{
			User: user,
		},
	}}

	s.registerClient(client)
	defer s.unregisterClient(client)

	closeChannel := make(chan error, 1)
	defer close(closeChannel)

	go client.Listen(closeChannel)

	err := <-closeChannel

	if err != nil {
		slog.Error("Closing connection due to error with client", "error", err)
		client.Close(err.Error())
	}

	client.Close("")
}

// registerClient registers a client
func (s *WebsocketService) registerClient(client *connection.ServerConnection) {
	connectedClientsMutex.Lock()

	s.connectedClients[client] = true

	defer connectedClientsMutex.Unlock()
}

// unregisterClient unregisters a client
func (s *WebsocketService) unregisterClient(client *connection.ServerConnection) {
	connectedClientsMutex.Lock()

	delete(s.connectedClients, client)

	defer connectedClientsMutex.Unlock()
}

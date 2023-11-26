package handlers

import (
	"fmt"
	"net/http"

	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	// Upgrader is used to upgrade a normal connection to a websocket
	upgrader websocket.Upgrader
	service  services.WebsocketService
}

// NewWebsocketHandler will instantiate a new WebsocketHandler
func NewWebsocketHandler(service services.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		service: service,
	}
}

// Establish will establish a websocket connection
func (h *WebsocketHandler) Establish(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Assert the user type
	userObject, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error trying to upgrade connection to websocket: %s", err).Error()})
		return
	}

	// Handle the WebSocket connection (e.g., register the connection, manage clients, etc.)
	go h.service.HandleConnection(conn, *userObject)
}

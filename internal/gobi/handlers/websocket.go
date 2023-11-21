package handlers

import "github.com/gin-gonic/gin"

type WebsocketHandler struct {
}

// NewWebsocketHandler will instantiate a new WebsocketHandler
func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}

// Establish will establish a websocket connection
func (h *WebsocketHandler) Establish(c *gin.Context) {

}

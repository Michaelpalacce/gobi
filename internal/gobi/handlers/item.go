package handlers

import (
	"net/http"

	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/gin-gonic/gin"
)

type ItemsHandler struct {
	Service *services.ItemsService
}

func NewItemsHandler(service *services.ItemsService) *ItemsHandler {
	return &ItemsHandler{
		Service: service,
	}
}

func (h *ItemsHandler) GetItem(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Not Implemented Yet"})
}

package handlers

import (
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/gin-gonic/gin"
)

type ItemsHandler struct {
	Service *services.ItemsService
}

// NewItemsHandler will instantiate a new ItemsHandler given the ItemService
func NewItemsHandler(service *services.ItemsService) *ItemsHandler {
	return &ItemsHandler{
		Service: service,
	}
}

func (h *ItemsHandler) AddItem(c *gin.Context) {

}

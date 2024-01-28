package handlers

import (
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/gin-gonic/gin"
)

// ItemHandler is the handler for the item routes
// @TODO: Implement
type ItemHandler struct {
	Service *services.ItemService
}

// NewItemHandler will instantiate a new ItemHandler given the ItemService
func NewItemHandler(service *services.ItemService) *ItemHandler {
	return &ItemHandler{
		Service: service,
	}
}

// GetItem will retrieve the item from the database
func (h *ItemHandler) GetItem(c *gin.Context) {
}

// CreateItem will insert the given item in the database.
// Returns 201 if the item is created successfully
func (h *ItemHandler) CreateItem(c *gin.Context) {
}

// DeleteItem will delete the item if it exists. If it does not exist, it will do nothing, but still return 200
func (h *ItemHandler) DeleteItem(c *gin.Context) {
}

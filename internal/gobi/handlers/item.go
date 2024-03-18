package handlers

import (
	"fmt"
	"net/http"

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
	form, _ := c.MultipartForm()
	files := form.File["item"]
	itemName := form.Value["itemName"]
	itemPath := form.Value["itemPath"]

	c.SaveUploadedFile(files[0], fmt.Sprintf("uploads/%s/%s", itemPath, itemName))
	c.String(http.StatusCreated, fmt.Sprintf("%d uploaded!", itemName))
}

// DeleteItem will delete the item if it exists. If it does not exist, it will do nothing, but still return 200
func (h *ItemHandler) DeleteItem(c *gin.Context) {
}

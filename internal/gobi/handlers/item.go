package handlers

import (
	"log/slog"
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
// @TODO: We need to make sure that when the user is uploading a file, we are saving it in the correct location and they have permission to do so
func (h *ItemHandler) CreateItem(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["item"]

	for _, file := range files {
		slog.Info("Uploading file", "filename", file.Filename)
		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			slog.Error("Error saving file", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully"})
}

// DeleteItem will delete the item if it exists. If it does not exist, it will do nothing, but still return 200
func (h *ItemHandler) DeleteItem(c *gin.Context) {
}

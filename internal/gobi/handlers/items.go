package handlers

import (
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
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

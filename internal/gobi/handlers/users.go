package handlers

import (
	"net/http"

	"github.com/Michaelpalacce/gobi/internal/gobi/models"
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	Service *services.UserService
}

// NewUsersHandler will instantiate a new UsersHandler given the UserService
func NewUsersHandler(service *services.UserService) *UsersHandler {
	return &UsersHandler{
		Service: service,
	}
}

// CreateUser will insert the given user in the database.
func (h UsersHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.Service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.Data(http.StatusCreated, "application/json", []byte{})
}

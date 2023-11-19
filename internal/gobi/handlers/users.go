package handlers

import (
	"fmt"
	"net/http"

	"github.com/Michaelpalacce/gobi/internal/gobi/models"
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	Service *services.UsersService
}

// NewUsersHandler will instantiate a new UsersHandler given the UserService
func NewUsersHandler(service *services.UsersService) *UsersHandler {
	return &UsersHandler{
		Service: service,
	}
}

// CreateUser will insert the given user in the database.
// Returns 201 if the user is created successfully
func (h UsersHandler) CreateUser(c *gin.Context) {
	var user = &models.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to bind user: %s", err).Error()})
		return
	}

	if err := h.Service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to create user: %s", err).Error()})
		return
	}

	c.Data(http.StatusCreated, "application/json", []byte{})
}

func (h UsersHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	h.Service.DeleteUser(id)
}

func (h UsersHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetUser(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to fetch user: %s", err).Error()})
	}

	c.JSON(http.StatusOK, user)
}

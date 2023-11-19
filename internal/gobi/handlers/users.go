package handlers

import (
	"fmt"
	"net/http"

	"github.com/Michaelpalacce/gobi/internal/gobi/models"
	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
// The password will be hashed, so we don't store it
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

	c.JSON(http.StatusCreated, bson.D{{Key: "_id", Value: user.ID}})
}

// DeleteUser will delete the user if it exists. If it does not exist, it will do nothing, but still return 200
// TODO: Only delete current user
// TODO: Delete user files and metadata in the future
func (h UsersHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to delete user: %s", err).Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte{})
}

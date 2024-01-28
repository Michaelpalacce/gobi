package handlers

import (
	"fmt"
	"net/http"

	"github.com/Michaelpalacce/gobi/internal/gobi/services"
	"github.com/Michaelpalacce/gobi/pkg/models"
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
	user := &models.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to bind user: %w", err).Error()})
		return
	}

	if err := h.Service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to create user: %w", err).Error()})
		return
	}

	c.JSON(http.StatusCreated, bson.D{{Key: "_id", Value: user.ID}})
}

// DeleteUser will delete the user if it exists. If it does not exist, it will do nothing, but still return 200
// TODO: Delete user files and metadata in the future
func (h UsersHandler) DeleteUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Assert the user type
	userObject, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if err := h.Service.DeleteUser(userObject.ID.Hex()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error while trying to delete user: %w", err).Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte{})
}

package services

import (
	"log/slog"

	"github.com/Michaelpalacce/gobi/internal/gobi/models"
	"github.com/Michaelpalacce/gobi/pkg/gobi/database"
)

type UsersService struct {
	DB *database.Database
}

// NewUsersService will return an instance of the User Service
func NewUsersService(db *database.Database) *UsersService {
	return &UsersService{
		DB: db,
	}
}

// CreateUser Create a new user in the database
// Will return an error if the user already exists or other issues happen
func (u UsersService) CreateUser(user models.User) error {
	slog.Info("Creating a new user", "user", user.Username)
	return nil
}

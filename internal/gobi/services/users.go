package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/internal/gobi/models"
	"github.com/Michaelpalacce/gobi/pkg/gobi/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
//
// # If a user is created successfully, then the InsertedID will be set in the passed user
//
// Will return an error if the user already exists or other issues happen
func (u UsersService) CreateUser(user *models.User) error {
	slog.Debug("Creating a new user", "user", user.Username)
	userCollection := u.DB.Collections.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := userCollection.FindOne(ctx, bson.D{{Key: "username", Value: user.Username}})

	if result.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("user exists")
	}

	insertResult, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error while inserting record: %s, error was %s", user.Username, err)
	}

	user.InsertedID = insertResult.InsertedID

	slog.Debug("User Created", "ID", user.InsertedID)

	return nil
}

// DeleteUser deletes the user given the ID.
// If the user doesn't exist, does nothing.
func (u UsersService) DeleteUser(id string) error {
	slog.Info("Deleting user", "id", id)
	return nil
}

// GetUser will return the user, given an ID.
// If the user does not exist, then an error will be returned.
func (u UsersService) GetUser(id string) (*models.User, error) {
	return &models.User{}, nil
}

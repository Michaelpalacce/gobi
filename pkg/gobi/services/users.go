package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/digest"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	slog.Info("Creating a new user", "user", user.Username)

	userCollection := u.DB.Collections.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := userCollection.FindOne(ctx, bson.D{{Key: "username", Value: user.Username}})

	if result.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("user exists")
	}

	// Set a new ObjectID for the user's _id field
	user.ID = primitive.NewObjectID()
	user.Password = digest.SHA256(user.Password)

	insertResult, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error while inserting record: %s, error was %w", user.Username, err)
	}

	id, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("error while type casting InsertedID")
	}

	user.ID = id

	slog.Info("User Created", "ID", user.ID)

	return nil
}

// DeleteUser deletes the user given the ID.
// If the user doesn't exist, does nothing.
func (u UsersService) DeleteUser(id string) error {
	slog.Info("Deleting user", "id", id)
	userCollection := u.DB.Collections.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: objectId}})

	return nil
}

// GetUser will return the user, given an ID.
// If the user does not exist, then an error will be returned.
func (u UsersService) GetUser(id string) (*models.User, error) {
	userCollection := u.DB.Collections.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: objectId}}).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, err
}

// GetUserByName will retrieve a user object given the username. Usernames should be unique
func (u UsersService) GetUserByName(username string) (*models.User, error) {
	userCollection := u.DB.Collections.UsersCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := &models.User{}

	err := userCollection.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

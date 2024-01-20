package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ItemsService handles the logic for CRUD operations on items
type ItemsService struct {
	DB *database.Database
}

func NewItemsService(db *database.Database) *ItemsService {
	return &ItemsService{
		DB: db,
	}
}

// Upsert will insert or update an item in the database
func (s *ItemsService) Upsert(item *models.Item, ownerId primitive.ObjectID) error {
	itemsCollection := s.DB.Collections.ItemCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := itemsCollection.FindOne(ctx, bson.D{
		{Key: "owner_id", Value: ownerId.Hex()},
		{Key: "server_path", Value: item.ServerPath},
		{Key: "vault_name", Value: item.VaultName},
	})

	doesNotExist := result.Err() == mongo.ErrNoDocuments

	if result.Err() != nil && !doesNotExist {
		return fmt.Errorf("error retrieving item: %v, error was %w", item, result.Err())
	}

	item.ServerMTime = time.Now().Unix()
	item.OwnerId = ownerId.Hex()
	// item.ID = primitive.NewObjectID()

	slog.Debug("Upserting item", "item", item)

	if !doesNotExist {
		resultRaw, err := result.Raw()
		if err != nil {
			return err
		}

		objectId := resultRaw.Lookup("_id").ObjectID()
		// item.ID = objectId

		return s.Update(objectId, item)
	}

	return s.Create(item)
}

// Create will insert an item into the database
func (s *ItemsService) Create(item *models.Item) error {
	itemsCollection := s.DB.Collections.ItemCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := itemsCollection.InsertOne(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

// Update will update an item in the database
func (s *ItemsService) Update(objectId primitive.ObjectID, item *models.Item) error {
	itemsCollection := s.DB.Collections.ItemCollection

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": item,
	}

	_, err := itemsCollection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		return err
	}

	return nil
}

package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database contains the mongo client and helpers to auth and execute queries
type Database struct {
	Client      *mongo.Client
	Collections collections

	Initialized  bool
	DatabaseName string
}

// NewDatabase represents a singleton instance of Database.
// It will initialize it and connect it to MongoDB
func NewDatabase() (*Database, error) {
	database := Database{
		Initialized: false,
	}

	if err := database.Connect(); err != nil {
		return nil, err
	}

	return &database, nil
}

// Connect validates that everything needed is present.
// It will also Ping the mongoDB instance
func (d *Database) Connect() error {
	if d.Initialized {
		return nil
	}

	if err := d.checkEnv(); err != nil {
		return err
	}

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return err
	}

	d.Client = client

	if err := d.Ping(); err != nil {
		return err
	}

	d.DatabaseName = os.Getenv("MONGO_DATABASE")
	d.Collections = newCollections(d)
	d.Initialized = true

	slog.Info("Connection to the Database was successful.")

	return nil
}

// Disconnects the MongoDB Database
func (d *Database) Disconnect() {
	slog.Info("Disconnecting from Database")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.Client.Disconnect(ctx); err != nil {
		panic(err)
	}

	d.Initialized = false
}

// Ping will send a ping request to the MongoDB Instance
func (d Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := d.Client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	slog.Debug("MongoDB pinged")

	return nil
}

// checkEnv will check if all the needed environment properties are set and return an error if they are not set
func (d Database) checkEnv() error {
	tokens := []string{"MONGO_CONNECTION_STRING", "MONGO_DATABASE"}

	for _, token := range tokens {
		_, exists := os.LookupEnv(token)

		if !exists {
			return fmt.Errorf("%s not set", token)
		}
	}

	return nil
}

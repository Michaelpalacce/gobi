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

// NewDatabase will create a new instance of Database and connecta to MongoDB
func NewDatabase() (*Database, error) {
	database := Database{
		Initialized: false,
	}

	if err := database.Init(); err != nil {
		return nil, err
	}

	return &database, nil
}

type Database struct {
	Client      *mongo.Client
	Initialized bool

	databaseName string
}

// Init validates that everything needed is present.
// It will also Ping the mongoDB instance
func (d *Database) Init() error {
	if d.Initialized {
		return nil
	}

	tokens := []string{"MONGO_CONNECTION_STRING", "MONGO_DATABASE"}

	for _, token := range tokens {
		_, exists := os.LookupEnv(token)

		if !exists {
			return fmt.Errorf("%s not set", token)
		}
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
	d.databaseName = os.Getenv("MONGO_DATABASE")

	if err := d.Ping(); err != nil {
		return err
	}

	d.Initialized = true

	slog.Info("Connection to the Database was successful.", "Database", d.databaseName)

	return nil
}

// Ping will send a ping request to the MongoDB Instance
func (d Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := d.Client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	slog.Debug("MongoDB pinged", "databaseName", d.databaseName)

	return nil
}

// Disconnects the MongoDB Database
func (d Database) Disconnect() {
	slog.Info("Disconnecting from Database", "databaseName", d.databaseName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.Client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

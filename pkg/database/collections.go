package database

import "go.mongodb.org/mongo-driver/mongo"

type collections struct {
	UsersCollection *mongo.Collection
	ItemCollection  *mongo.Collection
}

// newCollections will create a new Collections container that will contain all the possible collections supported by gobi
func newCollections(db *Database) collections {
	return collections{
		UsersCollection: db.Client.Database(db.DatabaseName).Collection("Users"),
		ItemCollection:  db.Client.Database(db.DatabaseName).Collection("Items"),
	}
}

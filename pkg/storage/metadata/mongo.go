package metadata

import (
	"context"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/database"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
)

// MongoDriver controls metadata for files in the MongoDB
type MongoDriver struct {
	DB     *database.Database
	Client *client.Client
}

// Reconcile will return all files that were changed since last sync
func (d *MongoDriver) Reconcile(lastSync int) ([]storage.Item, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"server_m_time": bson.M{"$gt": lastSync}},
			{"owner_id": d.Client.User.ID.Hex()},
			{"vault_name": d.Client.VaultName},
		},
	}
	// Perform the find operation
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cursor, err := d.DB.Collections.ItemCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the results
	var results []storage.Item
	for cursor.Next(context.Background()) {
		var result models.Item

		err := cursor.Decode(&result)

		if err != nil {
			return nil, err
		}

		results = append(results, storage.Item{Item: result})
	}

	// Check for errors from iterating over cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

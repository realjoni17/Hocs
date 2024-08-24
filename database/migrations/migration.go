package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RunMigration(db *mongo.Database) {
	// Example of creating indexes or ensuring collections exist
	collections := []string{"users", "services", "usercarts"}

	for _, collectionName := range collections {
		// This is where you might create indexes or initial data
		collection := db.Collection(collectionName)

		// Example: Creating an index (if needed)
		_, err := collection.Indexes().CreateOne(
			context.TODO(),
			mongo.IndexModel{
				Keys: bson.M{"field": 1}, // replace "field" with the actual field
			},
		)
		if err != nil {
			log.Fatalf("Failed to create index on collection %s: %v", collectionName, err)
		}

		log.Printf("Collection '%s' is set up.", collectionName)
	}
}

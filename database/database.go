package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func StartDatabase() {
	// MongoDB connection string for a local instance
	str := "mongodb://localhost:27017"

	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(str))
	if err != nil {
		fmt.Println("Could not connect to MongoDB")
		log.Fatal("Error:", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Could not ping MongoDB")
		log.Fatal("Error:", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Note: Migrations are not typically used in MongoDB the way they are in SQL databases.
	// If you have schema initialization or data seeding, you can do it here.
}

func CloseConn() error {
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDatabase() *mongo.Database {
	// Replace `yourDatabaseName` with the name of your MongoDB database
	return client.Database("Hdocs")
}

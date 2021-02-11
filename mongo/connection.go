package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", "user", "password", "localhost:27017"))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to init mongo client: %s", err.Error())
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Failed to connect mongo client: %s", err.Error())
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to execute a ping via mongo client: %s", err.Error())
	}

	log.Println("Successfully connected to MongoDB!")
}

package mongo

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/zneix/zniksbot/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	mongoPort, isPort := utils.GetEnv("MONGO_PORT", false)
	if !isPort {
		mongoPort = "27017"
	}
	mongoUser, _ := utils.GetEnv("MONGO_USER", true)
	mongoPassword, _ := utils.GetEnv("MONGO_PASSWORD", true)
	mongoAuthDb, isAuthDb := utils.GetEnv("MONGO_AUTHDB", false)
	if !isAuthDb {
		mongoAuthDb = "admin"
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s", mongoUser, url.QueryEscape(mongoPassword), "localhost:"+mongoPort, mongoAuthDb))
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

	return client
}

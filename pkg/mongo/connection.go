package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zneix/zneixbot-twitch/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection() *Connection {
	// Initialize context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize Client itself

	// Parse config values
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

	// Prepare mongo client's options
	uri := fmt.Sprintf("mongodb://%s:%s", "localhost", mongoPort)
	credentials := options.Credential{
		AuthSource: mongoAuthDb,
		Username:   mongoUser,
		Password:   mongoPassword,
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)

	// Actually connect to the database and test connection with a ping
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatalf("Failed to init mongo client: %s\n", err.Error())
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect mongo client: %s\n", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to execute a ping via mongo client: %s\n", err.Error())
	}

	log.Println("connected to MongoDB")

	return &Connection{
		client: client,
		ctx:    ctx,
	}
}

func (conn Connection) Disconnect() {
	conn.client.Disconnect(conn.ctx)
}

func (conn *Connection) Collection(name CollectionName) *mongo.Collection {
	return conn.client.Database(databaseName).Collection(string(name))
}

package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Connection struct {
	client *mongo.Client
	ctx    context.Context
}

// TODO: Add a config entry for databaseName, could be useful for easier testing
const databaseName = "zneixbot-twitch"

type CollectionName string

const (
	CollectionNameChannels = CollectionName("channels")
)

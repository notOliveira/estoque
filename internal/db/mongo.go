package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping
	if err := client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")

	return client, err
}

func GetCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}
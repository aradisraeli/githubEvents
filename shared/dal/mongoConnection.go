package dal

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectToMongo(mongoUri string) (*mongo.Client, error) {

	if mongoUri == "" {
		return nil, errors.New("you must set your 'MONGODB_URI' environment variable")
	}
	log.Println("Trying to connect to MongoDB.")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MongoDB.")
	return client, nil
}

func DisconnectFromMongo(client *mongo.Client) error {
	log.Println("Trying to disconnect from MongoDB.")
	if err := client.Disconnect(context.TODO()); err != nil {
		return err
	}
	log.Println("Successfully disconnected from MongoDB.")
	return nil
}

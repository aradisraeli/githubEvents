package dal

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoClient[T any] struct {
	mongo.Collection
}

func (x MongoClient[T]) ReplaceOneObj(object T, id string) error {
	filter := bson.D{{"_id", id}}
	opts := options.Replace().SetUpsert(true)
	result, err := x.ReplaceOne(context.TODO(), filter, object, opts)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 1 {
		log.Printf("Document with id %v was modified.", id)
	} else if result.UpsertedCount == 1 {
		log.Printf("Document with id %v was inserted.", id)
	}

	return nil
}

func (x MongoClient[T]) GetManyObj(skip int64, limit int64, filters map[string]any) ([]T, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)

	filtersBson := bson.M(filters)
	var results []T
	cursor, err := x.Find(context.TODO(), filtersBson, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	log.Printf("Got %v documents", len(results))
	return results, nil
}

func (x MongoClient[T]) GetOneObj(filters map[string]any) (*T, error) {
	filtersBson := bson.M(filters)
	var result T
	err := x.FindOne(context.TODO(), filtersBson).Decode(&result)
	if err != nil {
		return nil, err
	}
	log.Printf("Got one document: %v", result)
	return &result, nil
}

func (x MongoClient[T]) convertToInterfaces(objects []T) []interface{} {
	objectsInterface := make([]interface{}, len(objects))
	for i := range objects {
		objectsInterface[i] = objects[i]
	}
	return objectsInterface
}

func NewMongoClient[T any](client mongo.Client, dbName string, collectionName string) MongoClient[T] {
	collection := client.Database(dbName).Collection(collectionName)
	mongoClient := MongoClient[T]{*collection}
	return mongoClient
}

package controllers

import (
	"context"
	"githubEvents/shared"
	"githubEvents/shared/dal"
	"githubEvents/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllEvents(client *mongo.Client, page int64, size int64, params map[string]any) (models.Page[models.Event], error) {
	mongoEventsClient := dal.NewMongoClient[models.Event](*client, shared.MongoDatabaseName, shared.MongoEventsCollectionName)

	var paginatedResult models.Page[models.Event]
	paginatedResult.Page = page
	paginatedResult.Size = size
	totalEvents, err := mongoEventsClient.CountDocuments(context.TODO(), params)
	if err != nil {
		return models.Page[models.Event]{}, err
	}
	paginatedResult.Total = totalEvents
	events, err := mongoEventsClient.GetManyObj((page-1)*size, size, params)
	if err != nil {
		return models.Page[models.Event]{}, err
	}
	paginatedResult.Items = events
	paginatedResult.Count = len(events)

	return paginatedResult, nil
}

func CountEvents(client *mongo.Client, params map[string]any) (int64, error) {
	mongoEventsClient := dal.NewMongoClient[models.Event](*client, shared.MongoDatabaseName, shared.MongoEventsCollectionName)

	totalEvents, err := mongoEventsClient.CountDocuments(context.TODO(), params)
	if err != nil {
		return 0, err
	}
	return totalEvents, err
}

func RecentActors(client *mongo.Client, amount int64) ([]models.Actor, error) {
	mongoEventsClient := dal.NewMongoClient[models.Event](*client, shared.MongoDatabaseName, shared.MongoEventsCollectionName)

	groupStage := bson.D{{"$group", bson.D{
		{"_id", "$actor"},
		{"last_event_time", bson.D{{"$max", "$created_at"}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"last_event_time", -1}}}}
	limitStage := bson.D{{"$limit", amount}}
	projectStage := bson.D{{"$project", bson.D{
		{"id", "$_id.id"}, {"login", "$_id.login"}, {"display_login", "$_id.display_login"},
		{"gravatar_id", "$_id.gravatar_id"}, {"url", "$_id.url"}, {"avatar_url", "$_id.avatar_url"},
		{"_id", 0},
	}}}

	cursor, err := mongoEventsClient.Aggregate(context.TODO(), mongo.Pipeline{groupStage, sortStage, limitStage, projectStage})
	if err != nil {
		return nil, err
	}
	var results []models.Actor
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func RecentRepos(client *mongo.Client, amount int64) ([]models.EventRepo, error) {
	mongoEventsClient := dal.NewMongoClient[models.Event](*client, shared.MongoDatabaseName, shared.MongoEventsCollectionName)

	groupStage := bson.D{{"$group", bson.D{
		{"_id", "$repo"},
		{"last_event_time", bson.D{{"$max", "$created_at"}}},
		{"max_stars", bson.D{{"$max", "$repo.stars"}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"last_event_time", -1}}}}
	limitStage := bson.D{{"$limit", amount}}
	projectStage := bson.D{{"$project", bson.D{
		{"id", "$_id.id"}, {"name", "$_id.name"}, {"url", "$_id.url"},
		{"stars", "$max_stars"}, {"_id", 0},
	}}}

	cursor, err := mongoEventsClient.Aggregate(context.TODO(), mongo.Pipeline{groupStage, sortStage, limitStage, projectStage})
	if err != nil {
		return nil, err
	}
	var results []models.EventRepo
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

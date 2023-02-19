package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"taskexecutor/data"
	"taskexecutor/logger"
)

type Database struct {
	Database *mongo.Collection
}

func Connect(dbUrl string) (Database, error) {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return Database{}, err
	}

	taskQueueCollection := client.Database("tasks").Collection("task_queue")
	logger.Info("Database connection established")

	return Database{taskQueueCollection}, nil
}

func (db Database) InsertTask(ctx *context.Context, data *data.Task) error {
	_, err := db.Database.InsertOne(*ctx, data)
	if err != nil {
		return err
	}

	logger.Info("Inserting task with ID " + data.Id)
	return nil
}

func (db Database) UpdateTask(ctx *context.Context, data *data.Task) error {
	filter := bson.D{{"id", data.Id}}
	update := bson.D{{"$set", data}}

	result, err := db.Database.UpdateOne(*ctx, filter, update)

	if err != nil {
		return err
	}

	logger.Info("Modified tasks: " + fmt.Sprint(result.ModifiedCount))
	return nil
}

func (db Database) GetAllTasks(ctx *context.Context) ([]data.Task, error) {
	// Fetches all tasks as there is no filter criteria
	filter := bson.D{}
	cursor, err := db.Database.Find(*ctx, filter)

	if err != nil {
		return nil, err
	}

	var results []data.Task

	for cursor.Next(*ctx) {
		singleResult := data.Task{}

		err := cursor.Decode(&singleResult)
		if err != nil {
			return nil, err
		}

		results = append(results, singleResult)
	}

	if len(results) == 0 {
		logger.Warn("No tasks available")
	}

	return results, nil
}

func (db Database) QueryTaskByStatus(ctx *context.Context, status string, limit int) ([]data.Task, error) {
	filter := bson.M{"status": status}


	var filterOptions *options.FindOptions

	if limit != -1 {
		filterOptions = options.Find().SetLimit(int64(limit))
	}

	cursor, err := db.Database.Find(*ctx, filter, filterOptions)

	if err != nil {
		logger.Warn("Couldn't find any entry by status " + status)
	}

	var result []data.Task

	for cursor.Next(*ctx) {
		singleResult := data.Task{}

		err := cursor.Decode(&singleResult)
		if err != nil {
			return nil, err
		}

		result = append(result, singleResult)
	}

	return result, nil
}

func (db Database) GetTaskById(ctx *context.Context, id string) (data.Task, error) {
	filter := bson.M{"id": id}
	queryResult := db.Database.FindOne(*ctx, filter)

	if queryResult.Err() != nil {
		logger.Warn("Didn't find document with ID " + id)
		return data.Task{}, nil
	}

	var result data.Task
	err := queryResult.Decode(&result)

	if err != nil {
		return data.Task{}, err
	}

	return result, nil
}

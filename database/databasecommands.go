package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataWrapper[T any] struct {
	ID   primitive.ObjectID `bson:"_id"`
	Data T                  `bson:",inline"`
}

var EMPTY = mongo.ErrNoDocuments

type CollectionOptions func(col *mongo.Collection) error

func UniqueFields(fieldNames ...string) CollectionOptions {
	return func(col *mongo.Collection) error {
		if len(fieldNames) != 0 {
			keys := bson.D{}
			for _, fieldName := range fieldNames {
				keys = append(keys, primitive.E{Key: fieldName, Value: 1})
			}
			index, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: keys, Options: options.Index().SetUnique(true)})
			log.Printf("database index: %s created for collection: %s\n", index, col.Name())
			return err
		}
		return nil
	}
}

func NewCollection[T any](database DB, collection string, options ...CollectionOptions) (models CollectionCommands[T], events EventCommands[T]) {
	col := database.getDatabase().Collection(collection)
	for _, option := range options {
		err := option(col)
		if err != nil {
			log.Fatal(err)
		}
	}
	models = &mongoDBCollectionImpl[T]{
		database:   database.getDatabase(),
		collection: col,
	}
	events = &mongoDBEventsImpl[T]{
		database:   database.getDatabase(),
		collection: database.getDatabase().Collection(fmt.Sprintf("%s.event", collection)),
	}
	return
}

func insertOneIntoCollection[T any](object T, collection *mongo.Collection) (DataWrapper[T], error) {
	log.Printf("database performing InsertOne operation with: %+v\n", object)
	res, err := collection.InsertOne(context.TODO(), object)
	if err != nil {
		return DataWrapper[T]{}, err
	}
	return DataWrapper[T]{
		ID:   res.InsertedID.(primitive.ObjectID),
		Data: object,
	}, nil
}

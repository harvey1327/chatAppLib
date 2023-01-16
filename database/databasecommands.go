package database

import (
	"context"
	"fmt"
	"log"

	"github.com/harvey1327/chatapplib/models/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataWrapper[T any] struct {
	ID   primitive.ObjectID `bson:"_id"`
	Data T                  `bson:",inline"`
}

type CollectionCommands[T any] interface {
	FindByID(id string) (DataWrapper[T], error)
	FindByIDAndUpdate(object DataWrapper[T]) error
	FindSingleByQuery(query findBy) (DataWrapper[T], error)
	FindMultipleByQuery(query findBy) ([]DataWrapper[T], error)
	InsertOne(object T) (DataWrapper[T], error)
}

var EMPTY = mongo.ErrNoDocuments

type mongoDBCollectionImpl[T any] struct {
	database   *mongo.Database
	collection *mongo.Collection
}

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

func NewCollection[T any](database DB, collection string, options ...CollectionOptions) (models CollectionCommands[T], events CollectionCommands[message.EventMessage[T]]) {
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
	events = &mongoDBCollectionImpl[message.EventMessage[T]]{
		database:   database.getDatabase(),
		collection: database.getDatabase().Collection(fmt.Sprintf("%s.event", collection)),
	}
	return
}

func (m *mongoDBCollectionImpl[T]) InsertOne(object T) (DataWrapper[T], error) {
	log.Printf("database performing InsertOne operation with: %+v\n", object)
	res, err := m.collection.InsertOne(context.TODO(), object)
	if err != nil {
		return DataWrapper[T]{}, err
	}
	return DataWrapper[T]{
		ID:   res.InsertedID.(primitive.ObjectID),
		Data: object,
	}, nil
}

func (m *mongoDBCollectionImpl[T]) FindByID(id string) (DataWrapper[T], error) {
	log.Printf("database performing FindByID operation with: %s\n", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return DataWrapper[T]{}, err
	}
	var result DataWrapper[T]
	err = m.collection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		return DataWrapper[T]{}, err
	}
	return result, nil
}

type findBy map[string]interface{}

func Query(field string, value interface{}) findBy {
	m := make(findBy)
	m.And(field, value)
	return m
}

func (fb findBy) And(field string, value interface{}) {
	fb[field] = value
}

func (fb findBy) convert() bson.M {
	res := bson.M{}
	for k, v := range fb {
		res[k] = v
	}
	return res
}

func (m *mongoDBCollectionImpl[T]) FindSingleByQuery(query findBy) (DataWrapper[T], error) {
	log.Printf("database performing FindSingleByQuery operation with: %+v\n", query)
	var result DataWrapper[T]
	err := m.collection.FindOne(context.TODO(), query.convert()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return DataWrapper[T]{}, EMPTY
		}
		return DataWrapper[T]{}, err
	}
	return result, nil
}

func (m *mongoDBCollectionImpl[T]) FindByIDAndUpdate(object DataWrapper[T]) error {
	log.Printf("database performing FindByIDAndUpdate operation with: %+v\n", object)
	res := m.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": object.ID}, bson.M{"$set": object})
	return res.Err()
}

func (m *mongoDBCollectionImpl[T]) FindMultipleByQuery(query findBy) ([]DataWrapper[T], error) {
	log.Printf("database performing FindMultipleByQuery operation with: %+v\n", query)
	results := make([]DataWrapper[T], 0)
	curr, err := m.collection.Find(context.TODO(), query.convert())
	if err != nil {
		return nil, err
	}
	err = curr.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

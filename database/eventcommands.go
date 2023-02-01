package database

import (
	"context"
	"log"
	"time"

	"github.com/harvey1327/chatapplib/models/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventCommands[T any] interface {
	InsertOne(object message.EventMessage[T]) (DataWrapper[message.EventMessage[T]], error)
	ListenByEventID(id string) <-chan DataWrapper[message.EventMessage[T]]
}

type mongoDBEventsImpl[T any] struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func (m *mongoDBEventsImpl[T]) InsertOne(object message.EventMessage[T]) (DataWrapper[message.EventMessage[T]], error) {
	return insertOneIntoCollection(object, m.collection)
}

func (m *mongoDBEventsImpl[T]) ListenByEventID(id string) <-chan DataWrapper[message.EventMessage[T]] {
	log.Printf("database performing ListenByEventID operation with: %+v\n", id)
	events := make(chan DataWrapper[message.EventMessage[T]])

	go func() {
		cache := make(map[string]bool)
		for {
			results := make([]DataWrapper[message.EventMessage[T]], 0)
			curr, err := m.collection.Find(context.TODO(), bson.M{"eventID": id})
			if err != nil {
				panic(err)
			}
			err = curr.All(context.TODO(), &results)
			if err != nil {
				panic(err)
			}
			for _, res := range results {
				if _, ok := cache[res.ID.Hex()]; !ok {
					cache[res.ID.Hex()] = true
					events <- res
				}
			}
			if len(cache) == 2 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		close(events)
	}()
	return events
}

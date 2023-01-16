package message

import (
	"time"
)

type EventMessage[T any] struct {
	EventID   string    `json:"eventID" bson:"eventID"`
	Status    status    `json:"status" bson:"status"`
	Body      T         `json:"body,omitempty" bson:"body"`
	Error     string    `json:"error,omitempty" bson:"error"`
	ModelID   string    `json:"modelID,omitempty" bson:"modelID"`
	TimeStamp time.Time `json:"time" bson:"time"`
}

type status string

const (
	PENDING  status = "PENDING"
	COMPLETE status = "COMPLETE"
	FAILED   status = "FAILED"
)

func (event EventMessage[T]) Complete(modelID string) EventMessage[T] {
	return EventMessage[T]{
		EventID:   event.EventID,
		Status:    COMPLETE,
		ModelID:   modelID,
		TimeStamp: event.TimeStamp,
	}
}

func (event EventMessage[T]) Failed(reason string) EventMessage[T] {
	return EventMessage[T]{
		EventID:   event.EventID,
		Status:    FAILED,
		Body:      event.Body,
		Error:     reason,
		TimeStamp: event.TimeStamp,
	}
}

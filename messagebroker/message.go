package messagebroker

import (
	"time"
)

type EventMessage[T any] struct {
	EventID   string    `json:"eventID" bson:"eventID"`
	Status    status    `json:"status" bson:"status"`
	Body      T         `json:"body" bson:",inline"`
	Error     string    `json:"error,omitempty" bson:"error"`
	TimeStamp time.Time `json:"time" bson:"time"`
}

type status string

const (
	PENDING  status = "PENDING"
	COMPLETE status = "COMPLETE"
	FAILED   status = "FAILED"
)

func (event *EventMessage[T]) Complete() {
	event.Status = COMPLETE
}

func (event *EventMessage[T]) Failed(reason string) {
	event.Status = FAILED
	event.Error = reason
}

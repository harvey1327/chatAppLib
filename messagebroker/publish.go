package messagebroker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publish[T any] interface {
	Publish(message T) (EventMessage[T], error)
}

type rabbitPublish[T any] struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitPublisher[T any](broker MessageBroker, queueName string) Publish[T] {
	broker.declareQueue(queueName)
	return &rabbitPublish[T]{
		channel:   broker.getChannel(),
		queueName: queueName,
	}
}

func (rbtp *rabbitPublish[T]) Publish(message T) (EventMessage[T], error) {
	event := EventMessage[T]{Status: PENDING, Body: message, EventID: uuid.New().String(), TimeStamp: time.Now().UTC()}
	bytes, err := json.Marshal(event)
	if err != nil {
		return event, err
	}
	log.Printf("Publishing message %+v\n", event)
	return event, rbtp.channel.PublishWithContext(context.Background(), "", rbtp.queueName, false, false, amqp.Publishing{ContentType: "application/json", Body: bytes})
}

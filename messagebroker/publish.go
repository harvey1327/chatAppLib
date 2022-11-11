package messagebroker

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publish interface {
	Publish(message publishMessage) error
}

type rabbitPublish struct {
	channel *amqp.Channel
}

func NewRabbitPublish(broker MessageBroker) Publish {
	return &rabbitPublish{
		channel: broker.getChannel(),
	}
}

func (rbtp *rabbitPublish) Publish(message publishMessage) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	log.Printf("Publishing message %+v\n", message)
	return rbtp.channel.PublishWithContext(context.Background(), "", message.queueName, false, false, amqp.Publishing{ContentType: message.contentType, Body: bytes})
}

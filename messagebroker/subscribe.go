package messagebroker

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscribe[T any] interface {
	Subscribe(queueName string) <-chan EventMessage[T]
}

type rabbitSubscribe[T any] struct {
	channel *amqp.Channel
}

func NewRabbitSubscribe[T any](broker MessageBroker) Subscribe[T] {
	return &rabbitSubscribe[T]{
		channel: broker.getChannel(),
	}
}

func (rbtp *rabbitSubscribe[T]) Subscribe(queueName string) <-chan EventMessage[T] {
	log.Printf("Subscribing to %s\n", queueName)
	results := make(chan EventMessage[T])
	msgs, err := rbtp.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			received, ok := <-msgs
			if !ok {
				break
			}
			event := subscribeMessage[T]{}
			log.Printf("Read message: %+v\n", event)
			err := json.Unmarshal(received.Body, &event)
			if err != nil {
				log.Fatal(err)
			}
			results <- event.EventMessage
		}
		close(results)
	}()
	return results
}

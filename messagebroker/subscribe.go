package messagebroker

import (
	"encoding/json"
	"log"

	"github.com/harvey1327/chatapplib/models/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscribe[T any] interface {
	Subscribe() <-chan message.EventMessage[T]
}

type rabbitSubscribe[T any] struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitSubscriber[T any](broker MessageBroker, queueName string) Subscribe[T] {
	broker.declareQueue(queueName)
	return &rabbitSubscribe[T]{
		channel:   broker.getChannel(),
		queueName: queueName,
	}
}

func (rbtp *rabbitSubscribe[T]) Subscribe() <-chan message.EventMessage[T] {
	log.Printf("Subscribing to %s\n", rbtp.queueName)
	results := make(chan message.EventMessage[T])
	msgs, err := rbtp.channel.Consume(rbtp.queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			received, ok := <-msgs
			if !ok {
				break
			}
			event := message.EventMessage[T]{}
			log.Printf("Read message: %+v\n", event)
			err := json.Unmarshal(received.Body, &event)
			if err != nil {
				log.Fatal(err)
			}
			results <- event
		}
		close(results)
	}()
	return results
}

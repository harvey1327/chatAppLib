package messagebroker

import (
	"fmt"
	"log"

	"github.com/harvey1327/chatapplib/messagebroker/events/createuser"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	CloseConnection()
	getChannel() *amqp.Channel
}

type rabbitMessageBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

type messageBrokerConfig struct {
	host     string
	port     int
	username string
	password string
}

func MessageBrokerConfig(host string, port int, username string, password string) messageBrokerConfig {
	return messageBrokerConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (mbc messageBrokerConfig) validate() error {
	if mbc.host == "" {
		return fmt.Errorf("messagebroker host is invalid: '%s'", mbc.host)
	} else if mbc.port <= 0 {
		return fmt.Errorf("messagebroker port is invalid: '%d'", mbc.port)
	} else if mbc.username == "" {
		return fmt.Errorf("messagebroker username is invalid: '%s'", mbc.username)
	} else if mbc.password == "" {
		return fmt.Errorf("messagebroker password is invalid: '%s'", mbc.password)
	} else {
		return nil
	}
}

func NewRabbitMQ(config messageBrokerConfig) MessageBroker {
	err := config.validate()
	if err != nil {
		log.Fatal(err)
	}
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.username, config.password, config.host, config.port))
	if err != nil {
		log.Fatal(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	broker := &rabbitMessageBroker{
		connection: connection,
		channel:    channel,
	}

	broker.declareQueue(createuser.QUEUE_NAME)

	return broker
}

func (rmq *rabbitMessageBroker) declareQueue(queueName string) {
	_, err := rmq.channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (rmq *rabbitMessageBroker) CloseConnection() {
	rmq.channel.Close()
	rmq.connection.Close()
}

func (rmq *rabbitMessageBroker) getChannel() *amqp.Channel {
	return rmq.channel
}

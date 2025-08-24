package setup

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/env"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	rabbitAddress, err := env.GetEnv("AMQP_ADDRESS")
	if err != nil {
		return nil, nil, err
	}
	conn, err := amqp.Dial(rabbitAddress)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	return conn, ch, err
}

func DeclareExchange(ch *amqp.Channel, name, kind string) error {
	return ch.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

func DeclareQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(name, true, false, false, false, nil)
}

func BindQueue(ch *amqp.Channel, queue, routingKey, exchange string) error {
	return ch.QueueBind(queue, routingKey, exchange, false, nil)
}

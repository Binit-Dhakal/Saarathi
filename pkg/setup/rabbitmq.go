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

package setup

import amqp "github.com/rabbitmq/amqp091-go"

type QueueConfig struct {
	Name       string
	Exchange   string
	RoutingKey string
	Type       string
	Durable    bool
}

func SetupExchange(ch *amqp.Channel, name string, typeEx string, durable bool) error {
	err := ch.ExchangeDeclare(
		name,
		typeEx,
		durable,
		false, false, false, nil,
	)
	return err
}

func SetupQueues(ch *amqp.Channel, configs []QueueConfig) error {
	for _, cfg := range configs {
		if cfg.Exchange != "" {
			if err := SetupExchange(
				ch,
				cfg.Exchange,
				cfg.Type,
				cfg.Durable,
			); err != nil {
				return err
			}
		}

		queue, err := ch.QueueDeclare(cfg.Name, true, false, false, false, nil)
		if err != nil {
			return err
		}

		if cfg.Exchange != "" && cfg.RoutingKey != "" {
			if err := ch.QueueBind(queue.Name, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

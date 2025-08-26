package setup

import amqp "github.com/rabbitmq/amqp091-go"

type QueueConfig struct {
	Name       string
	Exchange   string
	RoutingKey string
	Type       string
	Durable    bool
}

func SetupQueues(ch *amqp.Channel, configs []QueueConfig) error {
	for _, cfg := range configs {
		if cfg.Exchange != "" {
			if err := ch.ExchangeDeclare(
				cfg.Exchange,
				cfg.Type,
				cfg.Durable,
				false, false, false, nil,
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

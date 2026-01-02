package rabbitmqmanager

import (
	"fmt"

	"cash-flow-financial/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQManager struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     *models.RabbitMQConfig
}

func NewRabbitMQManager(cfg *models.RabbitMQConfig) (*RabbitMQManager, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	rm := &RabbitMQManager{
		Connection: conn,
		Channel:    ch,
		Config:     cfg,
	}

	return rm, nil
}

func (rm *RabbitMQManager) Close() error {
	if rm.Channel != nil {
		rm.Channel.Close()
	}
	if rm.Connection != nil {
		rm.Connection.Close()
	}
	return nil
}

func (rm *RabbitMQManager) HealthCheck() error {
	if rm.Connection == nil || rm.Connection.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is closed")
	}
	return nil
}

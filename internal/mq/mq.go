package mq

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/retail-core/account-service/internal/logger"
	"go.uber.org/zap"
)

type RabbitMQConnection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func InitRabbitMQ(amqpURI string) (*RabbitMQConnection, error) {

	var conn *amqp091.Connection
	var err error

	config := amqp091.Config{
		Properties: amqp091.Table{
			"connection_name": "account-service",
		},
	}

	for i := 0; i < 10; i++ {

		conn, err = amqp091.DialConfig(amqpURI, config)
		if err == nil {
			break
		}

		logger.L().Error("Failed to connect to RabbitMQ, retrying...", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		logger.L().Fatal("Could not connect to RabbitMQ after several attempts", zap.Error(err))
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.L().Fatal("Failed to open a channel", zap.Error(err))
		return nil, err
	}

	logger.L().Info("Successfully connected to RabbitMQ", zap.String("amqpURI", amqpURI))
	return &RabbitMQConnection{conn: conn, channel: channel}, nil
}

func (r *RabbitMQConnection) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQConnection) Channel() *amqp091.Channel {
	return r.channel
}

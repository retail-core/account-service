package mq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp091.Channel
}

func NewPublisher(ch *amqp091.Channel) *Publisher {
	return &Publisher{ch: ch}
}

func (p *Publisher) PublishDomainEvent(ctx context.Context, routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		"domain.events", // ✅ publish to your domain exchange
		routingKey,      // e.g. "business_owner.created"
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) SetupDomainExchange() error {
	return p.ch.ExchangeDeclare(
		"domain.events", // exchange name
		"topic",         // type (allows pattern routing)
		true,            // durable
		false,           // auto-delete
		false,           // internal
		false,           // no-wait
		nil,
	)
}
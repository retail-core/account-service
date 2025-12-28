package consumer

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/retail-core/account-service/internal/dtos"
	"github.com/retail-core/account-service/internal/logger"
	"github.com/retail-core/account-service/internal/service"
	"go.uber.org/zap"
)

const (
	AccountExchangeName       = "domain.events"
	QueueName                 = "account_service_queue"
	BusinessCreatedRoutingKey = "business_owner.created"
	StaffCreatedRoutingKey    = "staff.created"
	StaffVerifiedRoutingKey   = "staff.verified"
)

type AccountConsumer struct {
	Ch      *amqp.Channel
	Service service.AccountService
}

func NewAccountConsumer(ch *amqp.Channel, service service.AccountService) *AccountConsumer {

	return &AccountConsumer{
		Ch:      ch,
		Service: service,
	}
}

func (c *AccountConsumer) StartConsumption(ctx context.Context) error {

	log := logger.L()

	// 1. Declare Exchange (ensures it exists)
	if err := c.Ch.ExchangeDeclare(AccountExchangeName, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	log.Info("📥 Declared Inventory exchange", zap.String("exchange Name", AccountExchangeName))

	q, err := c.Ch.QueueDeclare(QueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Info("📥 Declared Account Service queue", zap.String("queue Name", q.Name))

	// 3. Bind Queue to Exchange for confirmation and rollback keys
	if err = c.Ch.QueueBind(q.Name, BusinessCreatedRoutingKey, AccountExchangeName, false, nil); err != nil {
		return err
	}

	if err = c.Ch.QueueBind(q.Name, StaffCreatedRoutingKey, AccountExchangeName, false, nil); err != nil {
		return err
	}

	if err = c.Ch.QueueBind(q.Name, StaffVerifiedRoutingKey, AccountExchangeName, false, nil); err != nil {
		return err
	}

	msgs, err := c.Ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		log.Info("✅ Queue bound and consuming",
			zap.String("queue", q.Name),
			zap.String("routing_keys", BusinessCreatedRoutingKey+","+StaffCreatedRoutingKey+","+StaffVerifiedRoutingKey),
		)
		for m := range msgs {
			c.handleMessage(ctx, m)
		}
	}()

	<-ctx.Done()
	return nil
}

func (c *AccountConsumer) handleMessage(ctx context.Context, m amqp.Delivery) {

	var err error

	switch m.RoutingKey {
	case BusinessCreatedRoutingKey:
		var event dtos.BusinessCreatedEvent

		if err := json.Unmarshal(m.Body, &event); err != nil {
			m.Ack(false)
			return
		}

		var req dtos.CreateBusinessRequest = dtos.CreateBusinessRequest{
			UserID: event.OwnerId,
			Name:   event.OwnerName,
		}

		_, err = c.Service.CreateBusiness(ctx, req)
	case StaffCreatedRoutingKey:
		var event dtos.StaffCreatedEvent

		if err := json.Unmarshal(m.Body, &event); err != nil {
			m.Ack(false)
			return
		}

		var req dtos.CreateStaffRequest = dtos.CreateStaffRequest{
			AdminID:  event.OwnerID,
			UserID:   event.UserID,
			UserName: event.UserName,
			Email:    event.Email,
			Role:     &event.Role,
		}

		_, err = c.Service.CreateStaff(ctx, event.StoreID, req)

	case StaffVerifiedRoutingKey:
		var event dtos.VerifyStaffEvent

		if err := json.Unmarshal(m.Body, &event); err != nil {
			m.Ack(false)
			return
		}
		err = c.Service.VerifyStaff(ctx, event.UserID)

	default:
		m.Ack(false)
		logger.L().Warn(" [WARNING] Received message with unknown routing key: %s", zap.String("routingKey", m.RoutingKey))
		return
	}

	if err != nil {
		m.Nack(false, true)
		return
	}

	m.Ack(false)
	logger.L().Info(" [INFO] Successfully processed event %s", zap.String("routingKey", m.RoutingKey))
}

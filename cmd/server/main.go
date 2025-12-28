package main

import (
	"context"
	"net/http"

	"github.com/retail-core/account-service/internal/api"
	"github.com/retail-core/account-service/internal/config"
	"github.com/retail-core/account-service/internal/consumer"
	"github.com/retail-core/account-service/internal/db"
	"github.com/retail-core/account-service/internal/logger"
	"github.com/retail-core/account-service/internal/mq"
	"github.com/retail-core/account-service/internal/repository"
	"github.com/retail-core/account-service/internal/service"
	"go.uber.org/zap"
)

func main() {

	config := config.LoadConfig()

	logger.Init(config.Environment)

	_logger := logger.L()

	database, err := db.InitDB(config)
	if err != nil {
		_logger.Fatal("Database initialization failed", zap.Error(err))
	}

	rabbitConn, err := mq.InitRabbitMQ(config.RABBITMQ_URL)
	if err != nil {
		_logger.Fatal("Failed to initialize RabbitMQ connection", zap.Error(err))
	}

	defer rabbitConn.Close()

	r := api.ConfigureRoutes(database, rabbitConn)

	go func() {
		_logger.Info("sales service running on port", zap.String("PORT", config.Port))

		err = http.ListenAndServe(":"+config.Port, r)
		if err != nil {
			_logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	publisher := mq.NewPublisher(rabbitConn.Channel())
	publisher.SetupDomainExchange()

	accountRepo := repository.NewAccountRepository(database)
	accountService := service.NewAccountServiceImpl(accountRepo, publisher)

	consumer := consumer.NewAccountConsumer(rabbitConn.Channel(), accountService)
	consumer.StartConsumption(context.Background())

	select {}
}

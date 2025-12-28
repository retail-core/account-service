package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/retail-core/account-service/internal/handler"
	"github.com/retail-core/account-service/internal/mq"
	"github.com/retail-core/account-service/internal/repository"
	"github.com/retail-core/account-service/internal/service"
	"gorm.io/gorm"
)

func ConfigureRoutes(database *gorm.DB, rabbitConn *mq.RabbitMQConnection) http.Handler {

	accountRepo := repository.NewAccountRepository(database)

	publisher := mq.NewPublisher(rabbitConn.Channel())
	accountService := service.NewAccountServiceImpl(accountRepo, publisher)
	accountHandler := handler.NewAccountHandler(accountService)

	
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/v1", func(v1 chi.Router) {
		v1.Get("/users/{user_id}/stores", accountHandler.GetStoresByUserID)
		v1.Post("/users/{user_id}/stores", accountHandler.CreateStore)
		v1.Put("/users/{user_id}/stores/{store_id}", accountHandler.UpdateStore)
		v1.Post("/stores/{store_id}/staffs", accountHandler.CreateStaff)
		v1.Get("/stores/{store_id}/users", accountHandler.GetUsersByStoreID)
		v1.Get("/stores/{store_id}/staffs", accountHandler.GetStaffsByStoreID)
		v1.Post("/users/{user_id}/onboard", accountHandler.OnboardUser)
		v1.Delete("/stores/{store_id}/staffs/{staff_id}", accountHandler.DeleteStaff)
	})

	return r
}
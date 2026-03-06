package rest

import (
	"liquide_assignment/internal/service/authentication"
	"liquide_assignment/internal/service/login"
	"liquide_assignment/internal/service/order"
	"liquide_assignment/internal/service/report"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func InitHandlerNew(as authentication.AuthService, ls login.LoginService, os order.OrderService, rs report.ReportService) *chi.Mux {
	r := chi.NewRouter()

	allowedOrigins := []string{"http://localhost:8000"}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Mount("/v1", Routes(as, ls, os, rs))

	return r
}

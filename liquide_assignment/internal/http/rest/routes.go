package rest

import (
	"liquide_assignment/internal/service/authentication"
	"liquide_assignment/internal/service/login"
	"liquide_assignment/internal/service/order"
	"liquide_assignment/internal/service/report"

	"github.com/go-chi/chi"
)

func Routes(as authentication.AuthService, ls login.LoginService, os order.OrderService, rs report.ReportService) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/login", Login(ls))
		r.Post("/signup", Signup(ls))
	})

	r.Group(func(r chi.Router) {
		r.Use(as.ValidateSession)
		r.Post("/order", CreateOrder(os))
		r.Post("/orderbook", GetOrderBook(rs))
	})
	return r
}

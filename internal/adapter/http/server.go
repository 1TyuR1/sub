package http

import (
	"net/http"

	"crud_ef/internal/adapter/http/handlers"
	"crud_ef/internal/config"
	"crud_ef/internal/usecase/subscription"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	addr   string
	router *chi.Mux
}

func New(cfg config.Config, svc *subscription.Service) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	handlers.RegisterPing(r)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	sub := handlers.NewSubscriptionRoutes(svc)
	sub.Register(r)

	agg := handlers.NewAggregateRoutes(svc)
	agg.Register(r)

	return &Server{
		addr:   cfg.Addr(),
		router: r,
	}
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.addr, s.router)
}

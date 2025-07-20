package handler

import (
	"encoding/json"
	"net/http"
	"wallet-app/pkg/middlewares"
	"wallet-app/pkg/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middlewares.LoggingMiddleware)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallet", h.updateWalletBalance)
		r.Get("/wallets/{id}", h.getWalletInfo)
	})

	return r
}

// Standard Responses

func (h *Handler) sendError(w http.ResponseWriter, message string, status int) {
	if message == "" {
		http.Error(w, "", status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorRes{message})
}

func (h *Handler) sendSuccess(w http.ResponseWriter, message string, status int) {

	if message == "" {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(SuccessRes{message})
}

func (h *Handler) sendJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

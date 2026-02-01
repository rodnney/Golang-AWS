package handlers

import (
	"github.com/gorilla/mux"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

func SetupRouter(h *TransactionHandler, l logger.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(LoggingMiddleware(l))
	r.Use(RecoveryMiddleware(l))

	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/transactions", h.Create).Methods("POST")
	r.HandleFunc("/transactions/{id}", h.GetByID).Methods("GET")

	return r
}

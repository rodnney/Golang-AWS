package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/rodnney/transaction-processor/internal/services"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

type TransactionHandler struct {
	service *services.TransactionService
	logger  logger.Logger
}

func NewTransactionHandler(service *services.TransactionService, logger logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tx domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		h.logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTransaction(r.Context(), &tx); err != nil {
		h.logger.Error("Failed to create transaction", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tx, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get transaction", "id", id, "error", err)
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"OK"}`))
}

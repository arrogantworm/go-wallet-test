package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"wallet-app/pkg/models"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) getWalletInfo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	walletID, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, fmt.Sprintf("неверный id: %v", err), http.StatusBadRequest)
		return
	}
	balance := 1000.0
	res := models.Wallet{
		ID:      walletID,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
